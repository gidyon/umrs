package main

import (
	"context"
	"github.com/Sirupsen/logrus"
	"github.com/gidyon/config"
	patient_app "github.com/gidyon/umrs/internal/chaincodes/patient"
	contract_auth "github.com/gidyon/umrs/internal/pkg/ledger_contract"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/patient"
	"github.com/gidyon/micros"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func main() {
	// Root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// App config
	cfg, err := config.New()
	handleErr(err)

	// Create app
	app, err := micros.NewService(ctx, cfg)
	handleErr(err)

	app.SetServiceEndpoint("/")

	// Update app endpoints
	updateEndpoints(app)

	// Add htttp middlewares
	updateHTTPMiddlewares(app)

	// Adds interceptors and server options
	updateGRPCMiddlewares(app)

	// Initialize grpc server
	handleErr(app.InitGRPC(ctx))

	// Dial options
	contractID := uuid.New().String()
	contract := contract_auth.NewledgerContractAuth(contractID)
	dialOptions := []grpc.DialOption{
		grpc.WithPerRPCCredentials(contract),
		grpc.WithBlock(),
	}

	// Dial to ledger service
	cc, err := app.DialExternalService(ctx, "ledger", dialOptions)
	handleErr(err)

	logrus.Infoln("Connected to ledger")

	// patient chaincode
	patientChaincode, err := patient_app.NewPatientChaincode(ctx, &patient_app.Options{
		ContractID:       contractID,
		SQLDB:            app.GormDB(),
		ledgerClient: ledger.NewledgerClient(cc),
	})
	handleErr(err)

	// Register app modules to the gRPC server
	patient.RegisterPatientAPIServer(app.GRPCServer(), patientChaincode)

	// Register client(s) for reverse proxy server
	handleErr(patient.RegisterPatientAPIHandlerServer(ctx, app.RuntimeMux(), patientChaincode))

	// Run app
	handleErr(app.Run(ctx, false))
}

func handleErr(err error) {
	if err != nil {
		logrus.Fatalln(err)
	}
}