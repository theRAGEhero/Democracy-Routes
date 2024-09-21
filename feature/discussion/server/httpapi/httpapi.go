package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/httpapi/model"
)

func StartServer(tb testing.TB, port int) {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		tb.Helper()

		var auth model.UserAuthorizationResponse
		auth.Token = "authorized"

		require.NoError(tb, json.NewEncoder(w).Encode(auth), "encoding authorization response")
	})

	mux.HandleFunc("POST /meeting", func(w http.ResponseWriter, r *http.Request) {
		tb.Helper()

		var nm model.CreateMeeting

		require.NoError(tb, json.NewDecoder(r.Body).Decode(&nm), "decoding request")

		var m model.Meeting
		m.ID = "id"
		m.Name = nm.Name

		require.NoError(tb, json.NewEncoder(w).Encode(m), "encoding response")
	})

	var srv http.Server

	tb.Cleanup(func() {
		require.NoError(tb, srv.Shutdown(context.TODO()), "shutting down server")
	})

	srv.Addr = fmt.Sprintf("localhost:%d", port)

	srv.Handler = mux

	go srv.ListenAndServe()
}
