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

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/authenticationhandler"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli"
	createduser "github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli/command/root/create/user"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli/common"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/httpapi"
	apimodel "github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/httpapi/model"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/jwthandler"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/testhelper"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/userhandler"
	usermodel "github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/userhandler/model"
)

func TestServer(t *testing.T) {
	t.Parallel()

	t.Run("user authorization", func(t *testing.T) {
		t.Parallel()

		db := testhelper.TmpDB(t)

		authH, err := authenticationhandler.New(db)
		require.NoError(t, err, "creating auth handler")

		userH, err := userhandler.New(db)
		require.NoError(t, err, "creating user handler")

		jwtH, err := jwthandler.New([]byte("secret"))
		require.NoError(t, err, "creating jwt handler")

		api := httpApi(t, userH, authH, jwtH)

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
			Password: "secret",
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

		res, err = http.Get(api + "/")
		require.NoError(t, err, "getting client")
		require.Equal(t, http.StatusOK, res.StatusCode, "wrong status code")
		t.Cleanup(func() { require.NoError(t, res.Body.Close(), "closing response body") })
	})

	t.Run("new meeting", func(t *testing.T) {
		t.Parallel()

		t.Skip("TODO: implement")

		userH, err := userhandler.New(testhelper.TmpDB(t))
		require.NoError(t, err, "creating user handler")
		api := httpApi(t, nil, nil, nil)

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

func httpApi(tb testing.TB, userH *userhandler.Handler, authH *authenticationhandler.Handler, jwtH *jwthandler.Handler) string {
	tb.Helper()

	port := randomPort(tb)

	stop, err := httpapi.Start(httpapi.Settings{
		Port:            port,
		UserH:           userH,
		AuthenticationH: authH,
		JwtH:            jwtH,
	})
	require.NoError(tb, err, "starting http api")

	tb.Cleanup(func() {
		tb.Helper()

		require.NoError(tb, stop(context.TODO()), "stopping http api")
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
