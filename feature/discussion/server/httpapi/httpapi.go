package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/client"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/authenticationhandler"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/httpapi/model"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/jwthandler"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/meetinghandler"
	meetingmodel "github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/meetinghandler/model"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/userhandler"
)

type Settings struct {
	Port            int
	UserH           *userhandler.Handler
	AuthenticationH *authenticationhandler.Handler
	JwtH            *jwthandler.Handler
	MeetingH        *meetinghandler.Handler
}

func Start(settings Settings) (func(ctx context.Context) error, error) {
	if err := checkSettings(settings); err != nil {
		return nil, fmt.Errorf("checking settings: %w", err)
	}

	mux := http.NewServeMux()

	mux.Handle("GET /", http.FileServerFS(client.HTMLClient))

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		var req model.UserAuthorization

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpError(w, errors.New("decoding request"), http.StatusBadRequest)

			return
		}

		user, err := settings.UserH.GetByName(req.Username)
		if err != nil {
			httpError(w, errors.New("wrong credentials"), http.StatusUnauthorized)

			return
		}

		if !settings.AuthenticationH.Authenticate(user.ID, req.Password) {
			httpError(w, errors.New("wrong credentials"), http.StatusUnauthorized)

			return
		}

		token, err := settings.JwtH.Issue(user.ID)
		if err != nil {
			httpError(w, errors.New("issuing token error"), http.StatusInternalServerError)

			return
		}

		var auth model.UserAuthorizationResponse
		auth.Token = token

		if err := json.NewEncoder(w).Encode(auth); err != nil {
			httpError(w, fmt.Errorf("encoding authorization response: %w", err), http.StatusInternalServerError)

			return
		}
	})

	mux.HandleFunc("POST /meeting", func(w http.ResponseWriter, r *http.Request) {
		var nm model.CreateMeeting

		if err := json.NewDecoder(r.Body).Decode(&nm); err != nil {
			httpError(w, fmt.Errorf("decoding request: %w", err), http.StatusBadRequest)

			return
		}

		var m model.Meeting
		m.ID = "id"
		m.Name = nm.Title

		settings.MeetingH.Create(meetingmodel.CreateMeeting{
			Title: nm.Title,
		})

		if err := json.NewEncoder(w).Encode(m); err != nil {
			httpError(w, fmt.Errorf("encoding response: %w", err), http.StatusInternalServerError)

			return
		}
	})

	var srv http.Server

	srv.Addr = fmt.Sprintf("localhost:%d", settings.Port)

	srv.Handler = mux

	go srv.ListenAndServe()

	return srv.Shutdown, nil
}

func checkSettings(settings Settings) error {
	var aErr error

	if settings.Port == 0 {
		aErr = errors.Join(aErr, errors.New("port is not set"))
	}

	if settings.UserH == nil {
		aErr = errors.Join(aErr, errors.New("user handler is not set"))
	}

	if settings.AuthenticationH == nil {
		aErr = errors.Join(aErr, errors.New("authentication handler is not set"))
	}

	if settings.JwtH == nil {
		aErr = errors.Join(aErr, errors.New("jwt handler is not set"))
	}

	return aErr
}

type jsonError struct {
	Error string `json:"error"`
}

func httpError(w http.ResponseWriter, err error, code int) {
	h := w.Header()

	h.Del("Content-Length")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(jsonError{Error: err.Error()})
}
