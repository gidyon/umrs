package main

import (
	"context"
	"flag"
	"github.com/gidyon/config"
	"github.com/gidyon/micros"
)

var (
	insecure = flag.Bool("insecure", false, "Whether to use insecure http")
	mode     = flag.Bool("dev", false, "Run server in development mode")
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

	// Register app modules
	registerModules(ctx, app)

	// Run app
	handleErr(app.Run(ctx, *insecure))
}
