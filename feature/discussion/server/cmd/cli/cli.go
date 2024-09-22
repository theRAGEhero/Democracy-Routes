package main

import (
	"log"
	"os"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli"
)

func main() {
	if err := cli.Run(cli.Settings{
		Args: os.Args[1:],
		Out:  os.Stdout,
	}); err != nil {
		log.Fatalf("running app %s:", err)
	}
}
