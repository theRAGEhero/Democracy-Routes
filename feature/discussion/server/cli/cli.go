package cli

import (
	"fmt"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli/command/root"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli/common"
)

func Run(p common.Params) error {
	c, err := root.New(p)
	if err != nil {
		return fmt.Errorf("creating command: %w", err)
	}

	return c.Run()
}
