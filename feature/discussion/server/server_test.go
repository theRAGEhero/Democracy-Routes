package server_test

import (
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

		res, err := http.Post(srv, "text/html", strings.NewReader("authorization"))
		require.NoError(t, err, "authorizing user")
		require.Equal(t, http.StatusOK, res.StatusCode, "wrong status code")
		t.Cleanup(func() { require.NoError(t, res.Body.Close(), "closing response body") })

		// Then he can do it.

		var auth userAuthorizationResponse
		require.NoError(t, json.NewDecoder(res.Body).Decode(&auth), "decoding response body")

		assert.NotEmpty(t, auth.Token, "no authorization token")
	})
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

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tb.Helper()

		var auth userAuthorizationResponse
		auth.Token = "authorized"

		require.NoError(tb, json.NewEncoder(w).Encode(auth), "encoding authorization response")
	}))

	tb.Cleanup(srv.Close)

	return srv.URL
}
