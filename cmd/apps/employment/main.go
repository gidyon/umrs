package main

import (
	"context"
	"github.com/gidyon/config"
	employment_app "github.com/gidyon/umrs/internal/apps/employment"
	"github.com/gidyon/umrs/pkg/api/employment"
	"github.com/gidyon/micros"
)

func main() {
	// Root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// App config
	cfg, err := config.New()
	handleErr(err)

	// Create app utility
	app, err := micros.NewService(ctx, cfg, nil)
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

	employmentAPI, err := employment_app.NewEmploymentAPI(ctx, app.GormDB())
	handleErr(err)

	// Register app modules to the gRPC server
	employment.RegisterEmploymentAPIServer(app.GRPCServer(), employmentAPI)

	// Register client(s) for reverse proxy server
	handleErr(employment.RegisterEmploymentAPIHandlerServer(ctx, app.RuntimeMux(), employmentAPI))

	// Run app
	handleErr(app.Run(ctx))
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
