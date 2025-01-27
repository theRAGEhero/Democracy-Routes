package testhelper

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/openfga/go-sdk/client"
	"github.com/stretchr/testify/require"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/dbhandler"
)

const testContextTimeout = 10 * time.Second

func Context(tb testing.TB) context.Context {
	tb.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), testContextTimeout)
	tb.Cleanup(cancel)

	return ctx
}

func TmpDB(tb testing.TB) *sql.DB {
	tb.Helper()

	db := tmpDB(tb)
	prepareDB(tb, db)

	return db
}

func tmpDB(tb testing.TB) *sql.DB {
	tb.Helper()

	db, err := sql.Open("pgx", "postgres://postgres:testing@localhost:5432/postgres")
	require.NoError(tb, err, "connecting to db")
	require.NoError(tb, db.Ping(), "pinging db")

	tdbn := "db" + strings.ReplaceAll(uuid.NewString(), "-", "")

	_, err = db.Exec("CREATE DATABASE " + tdbn)
	require.NoError(tb, err, "creating temporary db")

	tdb, err := sql.Open("pgx", "postgres://postgres:testing@localhost:5432/"+tdbn)
	require.NoError(tb, err, "connecting to temporary db")
	require.NoError(tb, tdb.Ping(), "pinging temporary db")

	tb.Cleanup(func() {
		tb.Helper()

		require.NoError(tb, tdb.Close(), "closing temporary db connection")
		_, err = db.Exec("DROP DATABASE " + tdbn)
		require.NoError(tb, err, "dropping temporary db")
		require.NoError(tb, db.Close(), "closing db connection")
	})

	return tdb
}

func prepareDB(tb testing.TB, db *sql.DB) {
	tb.Helper()

	require.NoError(tb, dbhandler.PrepareDB(db), "preparing db")
}

func TmpOfgaClient(tb testing.TB) *client.OpenFgaClient {
	tb.Helper()

	ctx := Context(tb)

	const apiUrl = "http://localhost:6080"

	ofgaClient, err := client.NewSdkClient(&client.ClientConfiguration{
		ApiUrl: apiUrl,
	})
	require.NoError(tb, err, "creating a client")

	store, err := ofgaClient.CreateStore(ctx).Body(client.ClientCreateStoreRequest{Name: uuid.NewString()}).Execute()
	require.NoError(tb, err, "creating a store")

	err = ofgaClient.SetStoreId(store.Id)
	require.NoError(tb, err, "assigning store to the client")

	tb.Cleanup(func() {
		tb.Helper()

		_, err = ofgaClient.DeleteStore(ctx).Execute()
		require.NoError(tb, err, "deleting the store")
	})

	return ofgaClient
}
