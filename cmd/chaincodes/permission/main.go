package main

import (
	"context"
	"github.com/Sirupsen/logrus"
	"github.com/gidyon/config"
	permission_app "github.com/gidyon/umrs/internal/chaincodes/permission"
	contract_auth "github.com/gidyon/umrs/internal/pkg/ledger_contract"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/gidyon/umrs/pkg/api/permission"
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
	ledgerCC, err := app.DialExternalService(ctx, "ledger", dialOptions)
	handleErr(err)

	logrus.Infoln("Connected to ledger")

	// Dail to notification service
	notificationDialOptions := []grpc.DialOption{
		grpc.WithBlock(),
	}
	notificationCC, err := app.DialExternalService(ctx, "notification", notificationDialOptions)
	handleErr(err)

	logrus.Infoln("Connected to notification service")

	opt := &permission_app.Options{
		ContractID:         contractID,
		RedisClient:        app.RedisClient(),
		ledgerClient:   ledger.NewledgerClient(ledgerCC),
		NotificationClient: notification.NewNotificationServiceClient(notificationCC),
	}
	permissionAPI, err := permission_app.NewPermissionAPI(ctx, opt)
	handleErr(err)

	// Register app modules to the gRPC server
	permission.RegisterPatientPermissionAPIServer(app.GRPCServer(), permissionAPI)

	// Register client(s) for reverse proxy server
	handleErr(
		permission.RegisterPatientPermissionAPIHandlerServer(
			ctx, app.RuntimeMux(), permissionAPI,
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
