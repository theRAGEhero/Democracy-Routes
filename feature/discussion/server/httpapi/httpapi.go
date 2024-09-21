package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/httpapi/model"
)

func Start(port int) func(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		var auth model.UserAuthorizationResponse
		auth.Token = "authorized"

		if err := json.NewEncoder(w).Encode(auth); err != nil {
			http.Error(w, "encoding authorization response: "+err.Error(), http.StatusInternalServerError)

			return
		}
	})

	mux.HandleFunc("POST /meeting", func(w http.ResponseWriter, r *http.Request) {
		var nm model.CreateMeeting

		if err := json.NewDecoder(r.Body).Decode(&nm); err != nil {
			http.Error(w, "decoding request: "+err.Error(), http.StatusBadRequest)

			return
		}

		var m model.Meeting
		m.ID = "id"
		m.Name = nm.Name

		if err := json.NewEncoder(w).Encode(m); err != nil {
			http.Error(w, "encoding response: "+err.Error(), http.StatusInternalServerError)

			return
		}
	})

	var srv http.Server

	srv.Addr = fmt.Sprintf("localhost:%d", port)

	srv.Handler = mux

	go srv.ListenAndServe()

	return srv.Shutdown
}
