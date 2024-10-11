package main

import (
	"log"
	"os"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli/common"
)

func main() {
	if err := cli.Run(common.Params{
		Args: os.Args[1:],
		Out:  os.Stdout,
	}); err != nil {
		log.Fatalf("running app %s:", err)
	}
}
