package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatalf("running app %s:", err)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("cli", flag.ContinueOnError)

	err := fs.Parse(args)
	if err != nil {
		return fmt.Errorf("parsing flags: %w", err)
	}

	return nil
}
