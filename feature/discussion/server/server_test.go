package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		var auth userAuthorizationResponse
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

		var nm newMeeting
		nm.Name = "meeting"

		b, err := json.Marshal(nm)
		require.NoError(t, err, "marshalling request")

		res, err := http.Post(srv+"/meeting", "application/json", bytes.NewReader(b))
		require.NoError(t, err, "creating meeting")
		require.Equal(t, http.StatusOK, res.StatusCode, "wrong status code")
		t.Cleanup(func() { require.NoError(t, res.Body.Close(), "closing response body") })

		// Then he can do it.

		var m meeting
		require.NoError(t, json.NewDecoder(res.Body).Decode(&m), "decoding response")

		assert.NotEmpty(t, m.ID, "no meeting id")
		assert.Equal(t, nm.Name, m.Name, "wrong meeting name")
	})
}

type newMeeting struct {
	Name string
}

type meeting struct {
	ID   string
	Name string
}

type user struct{}

type userAuthorizationResponse struct {
	Token string
}

type userRepo struct{}

func (r *userRepo) GetUser(id string) (*user, error) {
	return &user{}, nil
}

func testServer(tb testing.TB) string {
	tb.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		tb.Helper()

		var auth userAuthorizationResponse
		auth.Token = "authorized"

		require.NoError(tb, json.NewEncoder(w).Encode(auth), "encoding authorization response")
	})

	mux.HandleFunc("POST /meeting", func(w http.ResponseWriter, r *http.Request) {
		tb.Helper()

		var nm newMeeting

		require.NoError(tb, json.NewDecoder(r.Body).Decode(&nm), "decoding request")

		var m meeting
		m.ID = "id"
		m.Name = nm.Name

		require.NoError(tb, json.NewEncoder(w).Encode(m), "encoding response")
	})

	srv := httptest.NewServer(mux)

	tb.Cleanup(srv.Close)

	return srv.URL
}
