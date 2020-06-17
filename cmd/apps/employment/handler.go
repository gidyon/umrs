package main

import (
	"github.com/gidyon/umrs/internal/pkg/token"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/micros"
	"github.com/gidyon/micros/utils/healthcheck"
	"net/http"
)

func updateEndpoints(app *micros.Service) {
	// Readiness health check
	app.AddEndpoint("/api/umrs/employments/readyq/", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      app,
		Type:         healthcheck.ProbeReadiness,
		AutoMigrator: func() error { return nil },
	}))

	// Liveness health check
	app.AddEndpoint("/api/umrs/employments/liveq/", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      app,
		Type:         healthcheck.ProbeLiveNess,
		AutoMigrator: func() error { return nil },
	}))

	// Token
	app.AddEndpointFunc("/api/umrs/employments/token/", func(w http.ResponseWriter, r *http.Request) {
		token.Write(w, int32(ledger.Actor_INSURANCE), r.URL.Query().Get("account_id"))
	})
}
