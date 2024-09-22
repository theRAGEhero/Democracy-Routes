package cli

import (
	"fmt"
	"io"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli/command/root"
)

type Settings struct {
	Args []string
	Out  io.Writer
}

func Run(s Settings) error {
	c, err := root.New(s.Args)
	if err != nil {
		return fmt.Errorf("creating command: %w", err)
	}

	return c.Run()
}
