package main

import (
	"context"
	"github.com/Sirupsen/logrus"
	"github.com/gidyon/config"
	account_app "github.com/gidyon/umrs/internal/apps/account"
	"github.com/gidyon/umrs/pkg/api/account"
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/gidyon/umrs/pkg/api/subscriber"
	"github.com/gidyon/micros"
	"google.golang.org/grpc"
)

func main() {
	// Root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// App config
	cfg, err := config.New()
	handleErr(err)

	cfg.DisableLogger()

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

	// Notification client
	dialOptions := []grpc.DialOption{
		grpc.WithBlock(),
	}
	notificationCC, err := app.DialExternalService(ctx, "notification", dialOptions)
	handleErr(err)

	accountAPI, err := account_app.NewAccountAPI(ctx, "", &account_app.Options{
		AppName:            "umrs NETWORK",
		SQLDB:              app.GormDB(),
		RedisDB:            app.RedisClient(),
		NotificationClient: notification.NewNotificationServiceClient(notificationCC),
		SubscriberClient:   subscriber.NewSubscriberAPIClient(notificationCC),
	})
	handleErr(err)

	logrus.Infoln("Connected to notification service")

	// Register app modules to the gRPC server
	account.RegisterAccountAPIServer(app.GRPCServer(), accountAPI)

	// Register client(s) for reverse proxy server
	handleErr(account.RegisterAccountAPIHandlerServer(ctx, app.RuntimeMux(), accountAPI))

	// Run app
	handleErr(app.Run(ctx, false))
}

func handleErr(err error) {
	if err != nil {
		logrus.Fatalln(err)
	}
}
