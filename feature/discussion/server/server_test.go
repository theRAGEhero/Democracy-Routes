package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/httpapi/model"
)

func TestServer(t *testing.T) {
	t.Parallel()

	t.Run("authorization", func(t *testing.T) {
		t.Parallel()

		var repo userRepo
		srv := testServer(t)

		// Given there is a user Dima.

		_, err := repo.GetUser("Dima")
		require.NoError(t, err, "getting user")

		// When Dima authorises.

		res, err := http.Post(srv+"/login", "text/plain", strings.NewReader("authorization"))
		require.NoError(t, err, "authorizing user")
		require.Equal(t, http.StatusOK, res.StatusCode, "wrong status code")
		t.Cleanup(func() { require.NoError(t, res.Body.Close(), "closing response body") })

		// Then he can do it.

		var auth model.UserAuthorizationResponse
		require.NoError(t, json.NewDecoder(res.Body).Decode(&auth), "decoding response body")

		assert.NotEmpty(t, auth.Token, "no authorization token")
	})

	t.Run("meeting", func(t *testing.T) {
		t.Parallel()

		var repo userRepo
		srv := testServer(t)

		// Given there is a user Dima.

		_, err := repo.GetUser("Dima")
		require.NoError(t, err, "getting user")

		// When Dima creates a new meeting.

		var nm model.CreateMeeting
		nm.Name = "meeting"

		b, err := json.Marshal(nm)
		require.NoError(t, err, "marshalling request")

		res, err := http.Post(srv+"/meeting", "application/json", bytes.NewReader(b))
		require.NoError(t, err, "creating meeting")
		require.Equal(t, http.StatusOK, res.StatusCode, "wrong status code")
		t.Cleanup(func() { require.NoError(t, res.Body.Close(), "closing response body") })

		// Then he can do it.

		var m model.Meeting
		require.NoError(t, json.NewDecoder(res.Body).Decode(&m), "decoding response")

		assert.NotEmpty(t, m.ID, "no meeting id")
		assert.Equal(t, nm.Name, m.Name, "wrong meeting name")
	})
}

type user struct{}

type userRepo struct{}

func (r *userRepo) GetUser(id string) (*user, error) {
	return &user{}, nil
}

func testServer(tb testing.TB) string {
	tb.Helper()

	port := randomPort(tb)

	startServer(tb, port)

	return fmt.Sprintf("http://localhost:%d", port)
}

func startServer(tb testing.TB, port int) {
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

func randomPort(tb testing.TB) int {
	tb.Helper()

	l, err := net.Listen(
		"tcp", // network
		"",    // address: The address is empty to force random port selection.
	)
	require.NoError(tb, err, "listening")

	require.NoError(tb, l.Close(), "closing listener")

	return l.Addr().(*net.TCPAddr).Port
}
