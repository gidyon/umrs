package main

import (
	"github.com/gidyon/micros"
	"github.com/gidyon/micros/utils/healthcheck"
)

func updateEndpoints(app *micros.Service) {
	// Readiness health check
	app.AddEndpoint("/api/umrs/notifications/readyq/", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      app,
		AutoMigrator: func() error { return nil },
		Type:         healthcheck.ProbeReadiness,
	}))

	// Liveness health check
	app.AddEndpoint("/api/umrs/notifications/liveq/", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      app,
		AutoMigrator: func() error { return nil },
		Type:         healthcheck.ProbeLiveNess,
	}))
}
