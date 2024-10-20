package server_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/meetinghandler"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/testhelper"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/userhandler"
)

func TestServer(t *testing.T) {
	t.Parallel()

	t.Run("user authorization", func(t *testing.T) {
		t.Parallel()

		db := testhelper.TmpDB(t)

		api := httpApi(t, httpApiSettings(t, randomPort(t), db))

		// Given there is a user Dima.

		var buf bytes.Buffer

		err := cli.Run(common.Params{
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
		defer res.Body.Close()

		// Then he can do it.

		var auth apimodel.UserAuthorizationResponse
		require.NoError(t, json.NewDecoder(res.Body).Decode(&auth), "decoding response body")

		assert.NotEmpty(t, auth.Token, "no authorization token")
	})

	t.Run("new meeting", func(t *testing.T) {
		t.Parallel()

		db := testhelper.TmpDB(t)

		api := httpApi(t, httpApiSettings(t, randomPort(t), db))

		// Given there is a user Dima.

		err := cli.Run(common.Params{
			Args: []string{"create", "user", "-name=Dima", "--pass=secret"},
			Out:  io.Discard,
			DB:   db,
		})
		require.NoError(t, err, "creating user")

		loginReq, err := json.Marshal(apimodel.UserAuthorization{
			Username: "Dima",
			Password: "secret",
		})
		require.NoError(t, err, "marshalling authorization request")

		loginRes, err := http.Post(api+"/login", "application/json", bytes.NewReader(loginReq))
		require.NoError(t, err, "making login request")
		require.Equal(t, http.StatusOK, loginRes.StatusCode, "wrong status code")
		defer loginRes.Body.Close()

		var auth apimodel.UserAuthorizationResponse
		require.NoError(t, json.NewDecoder(loginRes.Body).Decode(&auth), "decoding login response")

		// When Dima creates a new meeting.

		var nm apimodel.CreateMeeting
		nm.Title = "meeting"

		b, err := json.Marshal(nm)
		require.NoError(t, err, "marshalling request")

		meetingReq, err := http.NewRequest("POST", api+"/meeting", bytes.NewReader(b))
		require.NoError(t, err, "creating meeting request")

		meetingReq.Header.Set("Authorization", "Bearer "+auth.Token)

		var client http.Client

		res, err := client.Do(meetingReq)
		require.NoError(t, err, "making create meeting request")
		require.Equal(t, http.StatusOK, res.StatusCode, "wrong status code")
		defer res.Body.Close()

		// Then he can do it.

		var m apimodel.Meeting
		require.NoError(t, json.NewDecoder(res.Body).Decode(&m), "decoding response")

		assert.NotEmpty(t, m.ID, "no meeting id")
		assert.Equal(t, nm.Title, m.Title, "wrong meeting title")
	})
}

func httpApiSettings(tb testing.TB, port int, db *sql.DB) httpapi.Settings {
	tb.Helper()

	authH, err := authenticationhandler.New(db)
	require.NoError(tb, err, "creating auth handler")

	userH, err := userhandler.New(db)
	require.NoError(tb, err, "creating user handler")

	jwtH, err := jwthandler.New([]byte("secret"))
	require.NoError(tb, err, "creating jwt handler")

	meetingH, err := meetinghandler.New(db)
	require.NoError(tb, err, "creating meeting handler")

	return httpapi.Settings{
		Port:            port,
		UserH:           userH,
		AuthenticationH: authH,
		JwtH:            jwtH,
		MeetingH:        meetingH,
	}
}

func httpApi(tb testing.TB, settings httpapi.Settings) string {
	tb.Helper()

	stop, err := httpapi.Start(settings)
	require.NoError(tb, err, "starting http api")

	tb.Cleanup(func() {
		tb.Helper()

		require.NoError(tb, stop(context.TODO()), "stopping http api")
	})

	addr := fmt.Sprintf("http://localhost:%d", settings.Port)

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
