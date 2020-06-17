package main

import (
	"context"
	"github.com/Sirupsen/logrus"
	"github.com/gidyon/config"
	hospital_app "github.com/gidyon/umrs/internal/chaincodes/hospital"
	contract_auth "github.com/gidyon/umrs/internal/pkg/ledger_contract"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/hospital"
	"github.com/gidyon/umrs/pkg/api/notification"
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

	logrus.Infoln("dialing external services ...")

	// Dial options
	contractID := uuid.New().String()
	contract := contract_auth.NewledgerContractAuth(contractID)

	// Dial to ledger service
	dialOptions := []grpc.DialOption{
		grpc.WithPerRPCCredentials(contract),
		grpc.WithBlock(),
	}
	ledgerCC, err := app.DialExternalService(ctx, "ledger", dialOptions)
	handleErr(err)

	logrus.Infoln("connected to ledger service")

	// Dail to notification service
	notificationDialOptions := []grpc.DialOption{
		grpc.WithBlock(),
	}
	notificationCC, err := app.DialExternalService(ctx, "notification", notificationDialOptions)
	handleErr(err)

	logrus.Infoln("connected to notification service")

	// hospital chaincode
	hospitalChaincode, err := hospital_app.NewHospitalAPIServer(ctx, &hospital_app.Options{
		ContractID:         contractID,
		SQLDB:              app.GormDB(),
		ledgerClient:   ledger.NewledgerClient(ledgerCC),
		NotificationClient: notification.NewNotificationServiceClient(notificationCC),
	})
	handleErr(err)

	// Register app modules to the gRPC server
	hospital.RegisterHospitalAPIServer(app.GRPCServer(), hospitalChaincode)

	// Register client(s) for reverse proxy server
	handleErr(hospital.RegisterHospitalAPIHandlerServer(ctx, app.RuntimeMux(), hospitalChaincode))

	// Run app
	handleErr(app.Run(ctx, false))
}

func handleErr(err error) {
	if err != nil {
		logrus.Fatalln(err)
	}
}
