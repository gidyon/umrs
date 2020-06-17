package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/Sirupsen/logrus"
	"github.com/gidyon/umrs/internal/pkg/auth"
	contract_auth "github.com/gidyon/umrs/internal/pkg/ledger_contract"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/micros/pkg/conn"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

const certFile = "/home/gideon/go/src/github.com/gidyon/umrs/certs/ledger/cert.pem"

var createReqs = flag.Int("n", 10, "number of requests to create")

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Dial options
	contractID := uuid.New().String()
	contract := contract_auth.NewledgerContractAuth(contractID)

	logrus.Infoln("initializing dial options ...")
	dialOptions := []grpc.DialOption{
		grpc.WithPerRPCCredentials(contract),
		grpc.WithBlock(),
	}

	logrus.Infoln("dialing to server ...")
	cc, err := conn.DialService(ctx, &conn.GRPCDialOptions{
		ServerName:  setIfEmpty(os.Getenv("ledger_SERVER"), "localhost"),
		Address:     setIfEmpty(os.Getenv("ledger_ADDRESS"), "localhost:9090"),
		TLSCertFile: setIfEmpty(os.Getenv("TLS_CERT_FILE"), certFile),
		DialOptions: dialOptions,
	})
	handleErr(err)
	defer cc.Close()

	ledgerClient := ledger.NewledgerClient(cc)

	ctxReg := auth.AddSuperAdminMD(ctx)
	p, err := auth.AuthenticateSuperAdmin(ctxReg)
	handleErr(err)

	logrus.Infoln("registering client ...")
	_, err = ledgerClient.RegisterContract(ctxReg, &ledger.RegisterContractRequest{
		SuperAdminId: p.ID, ContractId: contractID,
	})
	handleErr(err)

	logrus.Infoln("sending RPCs ...")
	ctx = outgoingCtx(ctx)

	logrus.Infoln("creating block")
	addRes, err := createBlock(ctx, ledgerClient)
	handleErr(err)
	logrus.Infof("done creating block; hash is %s\n", addRes.Hash)

	time.Sleep(5 * time.Second)

	logrus.Infoln("getting block")
	blockRes, err := getBlock(ctx, addRes.Hash, ledgerClient)
	handleErr(err)
	logrus.Infof("done getting block; hash is %s\n", blockRes.Hash)

	logrus.Infoln("load testing started")
	loadTestAddingBlock(*createReqs, cc)
	logrus.Infoln("load testing done")
}

func handleErr(err error) {
	if err != nil {
		logrus.Fatalln(err)
	}
}

func setIfEmpty(val, def string) string {
	if strings.TrimSpace(val) == "" {
		return def
	}
	return val
}

func randomActor() ledger.Actor {
	index := rand.Intn(len(ledger.Actor_name) - 1)
	if index == 0 {
		index = 1
	}
	return ledger.Actor(index)
}

func randomOperation() ledger.Operation {
	index := rand.Intn(len(ledger.Operation_name) - 1)
	if index == 0 {
		index = 1
	}
	return ledger.Operation(index)
}

var organizations = []string{
	"ORG1", "ORG2", "ORG3", "ORG4", "ORG5",
}

func randomOrganization() string {
	return organizations[rand.Intn(len(organizations)-1)]
}

func newTransaction() *ledger.Transaction {
	tx := &ledger.Transaction{
		Operation: randomOperation(),
		Creator: &ledger.ActorPayload{
			Actor:         randomActor(),
			ActorId:       uuid.New().String(),
			ActorFullName: randomdata.SillyName(),
		},
		Patient: &ledger.ActorPayload{
			ActorId:       uuid.New().String(),
			ActorFullName: randomdata.SillyName(),
		},
		Organization: &ledger.ActorPayload{
			ActorId:       uuid.New().String(),
			ActorFullName: randomdata.Street() + " hospital",
		},
	}
	bs, err := proto.Marshal(tx)
	if err != nil {
		logrus.Fatalln(err)
	}

	tx.Details = bs
	return tx
}

func outgoingCtx(ctx context.Context) context.Context {
	tokenStr, err := auth.GenToken(ctx, nil, auth.HospitalGroup, 0)
	if err != nil {
		logrus.Fatalln(err)
	}
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", tokenStr))
	return metadata.NewIncomingContext(ctx, md)
}

func createBlock(
	ctx context.Context, client ledger.ledgerClient,
) (*ledger.AddBlockResponse, error) {
	callOptions := []grpc.CallOption{
		grpc.WaitForReady(true),
	}

	addRes, err := client.AddBlock(ctx, &ledger.AddBlockRequest{
		Transaction: newTransaction(),
	}, callOptions...)
	if err != nil {
		return nil, err
	}

	return addRes, nil
}

func getBlock(
	ctx context.Context, blockHash string, client ledger.ledgerClient,
) (*ledger.Block, error) {
	getRes, err := client.GetBlock(ctx, &ledger.GetBlockRequest{
		Hash: blockHash,
	}, grpc.WaitForReady(true))
	if err != nil {
		return nil, err
	}

	return getRes, nil
}

func loadTestAddingBlock(reqs int, cc *grpc.ClientConn) {
	wg := &sync.WaitGroup{}
	errChan := make(chan error, 0)
	ctx := outgoingCtx(context.Background())
	ledgerClient := ledger.NewledgerClient(cc)

	t1 := time.Now()

	for i := 0; i < reqs; i++ {
		wg.Add(1)
		tx := newTransaction()
		go func(tx *ledger.Transaction) {
			defer wg.Done()
			addRes, err := createBlock(ctx, ledgerClient)
			if err != nil {
				errChan <- err
				return
			}
			// logrus.Infof("Block added: hash value is %s\n", addRes.GetHash())
			_ = addRes
		}(tx)
	}

	go func() {
		wg.Wait()
		close(errChan)
		logrus.Infoln("almost up...")
	}()

	var errors int
	for range errChan {
		errors++
	}

	took := time.Since(t1).Seconds()

	logrus.Infof("Done: took %v seconds; %d Errors reported", took, errors)
}
