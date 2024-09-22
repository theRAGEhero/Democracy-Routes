package root

import (
	"fmt"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli/command/root/create"
)

type Command struct {
	cmd  string
	args []string
}

func New(args []string) (*Command, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no command specified")
	}

	return &Command{
		cmd:  args[0],
		args: args[1:],
	}, nil
}

func (c *Command) Run() error {
	switch c.cmd {
	case "create":
		cc, err := create.New(c.args)
		if err != nil {
			return err
		}

		return cc.Run()
	default:
		return fmt.Errorf("unknown command: %s", c.cmd)
	}
}
