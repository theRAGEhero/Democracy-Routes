package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
		require.NoError(t, err, "getting user Dima")

		// When Dima authorises.

		_, err = http.Post(srv, "text/html", strings.NewReader("authorization"))
		require.NoError(t, err, "authorizing Dima")

		// Then he can do it.

	})
}

type user struct{}

type userRepo struct{}

func (r *userRepo) GetUser(id string) (*user, error) {
	return &user{}, nil
}

func testServer(tb testing.TB) string {
	tb.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))

	tb.Cleanup(srv.Close)

	return srv.URL
}
