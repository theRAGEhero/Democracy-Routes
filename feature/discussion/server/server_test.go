package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/authandler"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli"
	createduser "github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli/command/root/create/user"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli/common"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/httpapi"
	apimodel "github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/httpapi/model"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/testhelper"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/userhandler"
	usermodel "github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/userhandler/model"
)

func TestServer(t *testing.T) {
	t.Parallel()

	t.Run("user authorization", func(t *testing.T) {
		t.Parallel()

		db := testhelper.TmpDB(t)

		authH, err := authandler.New(db)
		require.NoError(t, err, "creating auth handler")

		userH, err := userhandler.New(db)
		require.NoError(t, err, "creating user handler")

		api := httpApi(t, userH, authH)

		// Given there is a user Dima.

		var buf bytes.Buffer

		err = cli.Run(common.Params{
			Args: []string{"create", "user", "-name=Dima", "--pass=secret"},
			Out:  &buf,
			DB:   db,
		})
		require.NoError(t, err, "creating user")

		var addedUser createduser.Response
		require.NoError(t, json.Unmarshal(buf.Bytes(), &addedUser), "unmarshalling create user response")

		// When Dima authorises.

		req := apimodel.UserAuthorization{
			Username: addedUser.Name,
			Password: addedUser.Password,
		}

		b, err := json.Marshal(req)
		require.NoError(t, err, "marshalling authorization request")

		res, err := http.Post(api+"/login", "application/json", bytes.NewReader(b))
		require.NoError(t, err, "authorizing user")
		require.Equal(t, http.StatusOK, res.StatusCode, "wrong status code")
		t.Cleanup(func() { require.NoError(t, res.Body.Close(), "closing response body") })

		// Then he can do it.

		var auth apimodel.UserAuthorizationResponse
		require.NoError(t, json.NewDecoder(res.Body).Decode(&auth), "decoding response body")

		assert.NotEmpty(t, auth.Token, "no authorization token")
	})

	t.Run("new meeting", func(t *testing.T) {
		t.Parallel()

		userH, err := userhandler.New(testhelper.TmpDB(t))
		require.NoError(t, err, "creating user handler")
		api := httpApi(t, nil, nil)

		// Given there is a user Dima.

		u, err := userH.Create(&usermodel.CreateUser{
			Name: "Dima",
		})
		require.NoError(t, err, "creating user")

		_, err = userH.Get(u.ID)
		require.NoError(t, err, "getting user")

		// When Dima creates a new meeting.

		var nm apimodel.CreateMeeting
		nm.Name = "meeting"

		b, err := json.Marshal(nm)
		require.NoError(t, err, "marshalling request")

		res, err := http.Post(api+"/meeting", "application/json", bytes.NewReader(b))
		require.NoError(t, err, "creating meeting")
		require.Equal(t, http.StatusOK, res.StatusCode, "wrong status code")
		t.Cleanup(func() { require.NoError(t, res.Body.Close(), "closing response body") })

		// Then he can do it.

		var m apimodel.Meeting
		require.NoError(t, json.NewDecoder(res.Body).Decode(&m), "decoding response")

		assert.NotEmpty(t, m.ID, "no meeting id")
		assert.Equal(t, nm.Name, m.Name, "wrong meeting name")
	})
}

func httpApi(tb testing.TB, userH *userhandler.Handler, authH *authandler.Handler) string {
	tb.Helper()

	port := randomPort(tb)

	shutdown := httpapi.Start(port, userH, authH)

	tb.Cleanup(func() {
		tb.Helper()

		require.NoError(tb, shutdown(context.TODO()), "shutting down http api")
	})

	addr := fmt.Sprintf("http://localhost:%d", port)

	for i := 0; ; i++ {
		res, _ := http.Get(addr + "/health")

		if res != nil && res.StatusCode == http.StatusOK {
			break
		}

		if i > 30 {
			tb.Fatal("http api failed to start")
		}

		time.Sleep(100 * time.Millisecond)
	}

	return addr
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
