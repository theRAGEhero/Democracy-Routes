package testhelper

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/dbhandler"
)

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
