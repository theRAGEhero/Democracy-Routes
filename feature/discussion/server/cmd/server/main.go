package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"
	"time"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/authenticationhandler"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/dbhandler"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/httpapi"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/jwthandler"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/meetinghandler"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/userhandler"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("running server: %s", err)
	}
}

func run() error {
	db, err := dbhandler.DbConnection()
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

	err = dbhandler.PrepareDB(db)
	if err != nil {
		return fmt.Errorf("preparing db: %w", err)
	}

	userH, err := userhandler.New(db)
	if err != nil {
		return fmt.Errorf("creating user handler: %w", err)
	}

	authenticationH, err := authenticationhandler.New(db)
	if err != nil {
		return fmt.Errorf("creating authentication handler: %w", err)
	}

	jwtH, err := jwthandler.New([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return fmt.Errorf("creating jwt handler: %w", err)
	}

	meetingH, err := meetinghandler.New(db)
	if err != nil {
		return fmt.Errorf("creating meeting handler: %w", err)
	}

	stop, err := httpapi.Start(httpapi.Settings{
		Port:            8080,
		UserH:           userH,
		AuthenticationH: authenticationH,
		JwtH:            jwtH,
		MeetingH:        meetingH,
	})
	if err != nil {
		return fmt.Errorf("starting http api: %w", err)
	}

	gs.add(func(ctx context.Context) error {
		err := stop(ctx)
		if err != nil {
			return fmt.Errorf("stopping http api: %w", err)
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
