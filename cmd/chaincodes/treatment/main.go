package main

import (
	"context"
	"github.com/Sirupsen/logrus"
	"github.com/gidyon/config"
	treatment_app "github.com/gidyon/umrs/internal/chaincodes/treatment"
	contract_auth "github.com/gidyon/umrs/internal/pkg/ledger_contract"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/treatment"
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
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(contract),
	}

	// Dial to ledger service
	cc, err := app.DialExternalService(ctx, "ledger", dialOptions)
	handleErr(err)

	logrus.Infoln("Connected to ledger")

	opt := &treatment_app.Options{
		ContractID:       contractID,
		ledgerClient: ledger.NewledgerClient(cc),
	}

	treatmentAPI, err := treatment_app.NewTreatmentChaincode(ctx, opt)
	handleErr(err)

	// Register app modules to the gRPC server
	treatment.RegisterTreatmentAPIServer(app.GRPCServer(), treatmentAPI)

	// Register client(s) for reverse proxy server
	handleErr(
		treatment.RegisterTreatmentAPIHandlerServer(
			ctx, app.RuntimeMux(), treatmentAPI,
		),
	)

	// Run app
	handleErr(app.Run(ctx, false))
}

func handleErr(err error) {
	if err != nil {
		logrus.Fatalln(err)
	}
}
