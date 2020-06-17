package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/Sirupsen/logrus"
	ledger_app "github.com/gidyon/umrs/internal/ledger"
	"github.com/gidyon/umrs/internal/pkg/md"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/micros/pkg/conn"
	app_grpc "github.com/gidyon/micros/pkg/grpc"
	app_grpc_middleware "github.com/gidyon/micros/pkg/grpc/middleware"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"os"
	"strings"
	// sqlite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	port     = flag.Int("p", 9090, "port for the ledger app")
	dbFile   = flag.String("f", "./db", "path to sqlite3 db file")
	useMysql = flag.Bool("mysql", false, "use mysql database to store ledger")
)

const (
	certFile = "/home/gideon/go/src/github.com/gidyon/umrs/certs/ledger/cert.pem"
	keyFile  = "/home/gideon/go/src/github.com/gidyon/umrs/certs/ledger/key.pem"
)

func main() {
	flag.Parse()

	ctx := context.Background()

	unaryInterceptors := make([]grpc.UnaryServerInterceptor, 0)
	streamInterceptors := make([]grpc.StreamServerInterceptor, 0)

	// Unary interceptor for authentication
	contractAuth := grpc.UnaryServerInterceptor(contractAuthInterceptor)
	unaryInterceptors = append(unaryInterceptors, contractAuth)

	// Recovery middleware
	recoveryUIs, recoverySIs := app_grpc_middleware.AddRecovery()
	unaryInterceptors = append(unaryInterceptors, recoveryUIs...)
	streamInterceptors = append(streamInterceptors, recoverySIs...)

	transportCreds, err := credentials.NewServerTLSFromFile(
		setIfEmpty(os.Getenv("TLS_CERT_FILE"), certFile), setIfEmpty(os.Getenv("TLS_KEY_FILE"), keyFile),
	)
	handleErr(err)

	serverOptions := []grpc.ServerOption{
		grpc.Creds(transportCreds),
	}

	s, err := app_grpc.NewGRPCServer(&app_grpc.ServerParams{
		ServerOptions:     serverOptions,
		UnaryInterceptors: unaryInterceptors,
		StreamInterceptos: streamInterceptors,
	})
	handleErr(err)

	var sqlDB *gorm.DB
	if !*useMysql {
		sqlDB, err = gorm.Open("sqlite3", setIfEmpty(os.Getenv("DB_FILE"), *dbFile))
		handleErr(err)
		logrus.Infoln("sqlite initialized")
	} else {
		sqlDB, err = conn.ToSQLDBUsingORM(&conn.DBOptions{
			Dialect:  "mysql",
			Host:     setIfEmpty(os.Getenv("MYSQL_HOST"), "localhost"),
			Port:     setIfEmpty(os.Getenv("MYSQL_PORT"), "3306"),
			User:     setIfEmpty(os.Getenv("MYSQL_USER"), "root"),
			Password: setIfEmpty(os.Getenv("MYSQL_PASSWORD"), "hakty11"),
			Schema:   setIfEmpty(os.Getenv("MYSQL_SCHEMA"), "umrs"),
		})
		handleErr(err)
		logrus.Infoln("mysql initialized")
	}

	opt := &ledger_app.Options{
		SQLDB: sqlDB,
		RedisClient: redis.NewClient(&redis.Options{
			Addr: setIfEmpty(os.Getenv("ORDERER_ADDRESS"), "localhost:6379"),
		}),
		Network:      setIfEmpty(os.Getenv("ledger_NETWORK"), "umrs_net"),
		List:         setIfEmpty(os.Getenv("ledger_LIST"), "umrs_list"),
		SymmetricKey: []byte(setIfEmpty(os.Getenv("SECRET_KEY"), randomdata.RandStringRunes(32))),
	}

	srv, err := ledger_app.NewledgerServer(ctx, opt)
	handleErr(err)

	ledger.RegisterledgerServer(s, srv)

	port := setIfEmpty(os.Getenv("SERVICE_PORT"), "9090")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", strings.TrimSuffix(port, ":")))
	handleErr(err)

	logrus.Infof("<gRPC server started on port %s>", port)

	handleErr(s.Serve(lis))
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

func contractAuthInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	if strings.Contains(info.FullMethod, "RegisterContract") {
		return handler(md.AddFromCtx(ctx), req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.NotFound, "metadata not found")
	}

	contractID := strings.Join(md["contract_id"], "")
	err := ledger_app.AuthorizeContract(contractID)
	if err != nil {
		return nil, err
	}

	ctx = metadata.NewOutgoingContext(ctx, md)
	return handler(ctx, req)
}
