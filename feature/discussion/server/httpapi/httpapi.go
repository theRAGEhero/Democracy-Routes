package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/authenticationhandler"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/httpapi/model"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/jwthandler"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/userhandler"
)

type Settings struct {
	Port  int
	UserH *userhandler.Handler
	AuthH *authenticationhandler.Handler
	JwtH  *jwthandler.Handler
}

func Start(settings Settings) func(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		var req model.UserAuthorization

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "decoding request", http.StatusBadRequest)

			return
		}

		user, err := settings.UserH.GetByName(req.Username)
		if err != nil {
			http.Error(w, "wrong credentials", http.StatusUnauthorized)

			return
		}

		if !settings.AuthH.Authenticate(user.ID, req.Password) {
			http.Error(w, "wrong credentials", http.StatusUnauthorized)

			return
		}

		token, err := settings.JwtH.Issue(user.ID)
		if err != nil {
			http.Error(w, "issuing token error", http.StatusInternalServerError)
		}

		var auth model.UserAuthorizationResponse
		auth.Token = token

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

	srv.Addr = fmt.Sprintf("localhost:%d", settings.Port)

	srv.Handler = mux

	go srv.ListenAndServe()

	return srv.Shutdown
}
