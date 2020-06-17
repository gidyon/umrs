package main

import (
	"encoding/json"
	"github.com/gidyon/gateway"
	"github.com/gidyon/umrs/internal/pkg/token"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"net/http"
)

func updateEndpoints(g *gateway.Gateway) {
	// Update health checks endpoints
	g.HandleFunc("/readyq", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := "am ready and running :)"
		w.Write([]byte(data))
	}))
	g.HandleFunc("/liveq", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := "am alive and running :)"
		w.Write([]byte(data))
	}))

	// Account groups
	g.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
		groups := []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{
			{Name: "Patient Account", Value: ledger.Actor_PATIENT.String()},
			{Name: "Hospital Account", Value: ledger.Actor_HOSPITAL.String()},
			{Name: "Insurance Account", Value: ledger.Actor_INSURANCE.String()},
			{Name: "Governement Account", Value: ledger.Actor_GOVERNMENT.String()},
			{Name: "Admin Account", Value: ledger.Actor_ADMIN.String()},
		}
		err := json.NewEncoder(w).Encode(groups)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Default account token handler
	g.HandleFunc("/token/default", func(w http.ResponseWriter, r *http.Request) {
		token.Write(w, int32(ledger.Actor_PATIENT), r.URL.Query().Get("actor_id"))
	})
	// Patient token
	g.HandleFunc("/token/patient", func(w http.ResponseWriter, r *http.Request) {
		token.Write(w, int32(ledger.Actor_PATIENT), r.URL.Query().Get("actor_id"))
	})
	// Insurance token
	g.HandleFunc("/token/insurance", func(w http.ResponseWriter, r *http.Request) {
		token.Write(w, int32(ledger.Actor_INSURANCE), r.URL.Query().Get("actor_id"))
	})
	// Patient token
	g.HandleFunc("/token/hospital", func(w http.ResponseWriter, r *http.Request) {
		token.Write(w, int32(ledger.Actor_HOSPITAL), r.URL.Query().Get("actor_id"))
	})
	// Patient token
	g.HandleFunc("/token/government", func(w http.ResponseWriter, r *http.Request) {
		token.Write(w, int32(ledger.Actor_GOVERNMENT), r.URL.Query().Get("actor_id"))
	})
	// Patient token
	g.HandleFunc("/token/admin", func(w http.ResponseWriter, r *http.Request) {
		token.Write(w, int32(ledger.Actor_ADMIN), r.URL.Query().Get("actor_id"))
	})
	g.Handle("/hello", http.HandlerFunc(helloHandler))

}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
