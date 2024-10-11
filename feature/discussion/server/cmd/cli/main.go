package main

import (
	"fmt"
	"log"
	"os"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli/common"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/dbhandler"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	db, err := dbhandler.DbConnection()
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}

	err = dbhandler.PrepareDB(db)
	if err != nil {
		return fmt.Errorf("preparing db: %w", err)
	}

	err = cli.Run(common.Params{
		Args: os.Args[1:],
		Out:  os.Stdout,
		DB:   db,
	})
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
