package main

import (
	"github.com/gidyon/logger"
	"github.com/gidyon/micros"
	app_grpc_middleware "github.com/gidyon/micros/pkg/grpc/middleware"
	http_middleware "github.com/gidyon/micros/pkg/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func updateGRPCMiddlewares(app *micros.Service) {
	// Logging middleware
	if app.Config().Logging() {
		loggingUIs, loggingSIs := app_grpc_middleware.AddLogging(logger.Log)
		app.AddGRPCUnaryServerInterceptors(loggingUIs...)
		app.AddGRPCStreamServerInterceptors(loggingSIs...)
	}

	// Recovery middleware
	recoveryUIs, recoverySIs := app_grpc_middleware.AddRecovery()
	app.AddGRPCUnaryServerInterceptors(recoveryUIs...)
	app.AddGRPCStreamServerInterceptors(recoverySIs...)

	transportCreds, err := credentials.NewServerTLSFromFile(
		app.Config().ServiceTLSCertFile(), app.Config().ServiceTLSKeyFile(),
	)
	handleErr(err)

	serverOptions := []grpc.ServerOption{
		grpc.Creds(transportCreds),
	}

	// Add server options
	app.AddGRPCServerOptions(serverOptions...)
}

func updateHTTPMiddlewares(app *micros.Service) {
	// Enable CORS
	if *mode {
		app.AddHTTPMiddlewares(http_middleware.SupportCORS)
	}
	app.AddHTTPMiddlewares(http_middleware.AddRequestID)
}
