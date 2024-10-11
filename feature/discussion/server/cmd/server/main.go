package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/dbhandler"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/httpapi"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/userhandler"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("running server: %s", err)
	}
}

func run() error {
	db, err := dbConnection()
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}

	var gs gracefulShutdown

	gs.add(func(ctx context.Context) error {
		err := db.Close()
		if err != nil {
			return fmt.Errorf("closing db: %w", err)
		}
		return nil
	})

	userH, err := userhandler.New(db)
	if err != nil {
		return fmt.Errorf("creating user handler: %w", err)
	}

	stop := httpapi.Start(httpapi.Settings{
		Port:  8080,
		UserH: userH,
		AuthH: nil,
		JwtH:  nil,
	})

	gs.add(func(ctx context.Context) error {
		err := stop(ctx)
		if err != nil {
			return fmt.Errorf("stopping server: %w", err)
		}
		return nil
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return gs.shutdown(ctx)
}

func dbConnection() (*sql.DB, error) {
	db, err := sql.Open("pgx", "postgres://postgres:testing@localhost:5432/postgres")
	if err != nil {
		return nil, fmt.Errorf("connecting to default postgres db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("pinging postgres db: %w", err)
	}

	const dbName = "democracy_routes"

	db.Exec("CREATE DATABASE " + dbName)

	err = db.Close()
	if err != nil {
		return nil, fmt.Errorf("closing postgres db: %w", err)
	}

	db, err = sql.Open("pgx", "postgres://postgres:testing@localhost:5432/"+dbName)
	if err != nil {
		return nil, fmt.Errorf("connecting to %s db: %w", dbName, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("pinging %s db: %w", dbName, err)
	}

	err = dbhandler.PrepareDB(db)
	if err != nil {
		return nil, fmt.Errorf("preparing %s db: %w", dbName, err)
	}

	return db, nil
}

type gracefulShutdown struct {
	actions []func(context.Context) error
}

func (g *gracefulShutdown) add(action func(context.Context) error) {
	g.actions = append(g.actions, action)
}

func (g *gracefulShutdown) shutdown(ctx context.Context) error {
	var aErr error

	slices.Reverse(g.actions)

	for _, action := range g.actions {
		if err := action(ctx); err != nil {
			aErr = errors.Join(aErr, err)
		}
	}

	return aErr
}
