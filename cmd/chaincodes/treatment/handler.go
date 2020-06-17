package main

import (
	"github.com/gidyon/micros"
	"github.com/gidyon/micros/utils/healthcheck"
)

func updateEndpoints(app *micros.Service) {
	// Readiness health check
	app.AddEndpoint("/api/umrs/treatments/readyq/", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      app,
		Type:         healthcheck.ProbeReadiness,
		AutoMigrator: func() error { return nil },
	}))

	// Liveness health check
	app.AddEndpoint("/api/umrs/treatments/liveq/", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      app,
		Type:         healthcheck.ProbeLiveNess,
		AutoMigrator: func() error { return nil },
	}))
}
