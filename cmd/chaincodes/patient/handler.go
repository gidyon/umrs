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
	app.AddEndpoint("/api/umrs/patients/readyq/", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      app,
		Type:         healthcheck.ProbeReadiness,
		AutoMigrator: func() error { return nil },
	}))

	// Liveness health check
	app.AddEndpoint("/api/umrs/patients/liveq/", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      app,
		Type:         healthcheck.ProbeLiveNess,
		AutoMigrator: func() error { return nil },
	}))

	// Token Generator for testing
	app.AddEndpointFunc("/api/umrs/patients/token/", func(w http.ResponseWriter, r *http.Request) {
		token.Write(w, int32(ledger.Actor_PATIENT), r.URL.Query().Get("patient_id"))
	})
	app.AddEndpointFunc("/api/umrs/patients/token/hospital", func(w http.ResponseWriter, r *http.Request) {
		token.Write(w, int32(ledger.Actor_HOSPITAL), r.URL.Query().Get("hospital_id"))
	})
}
