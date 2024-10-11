package create

import (
	"fmt"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli/command"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli/command/root/create/user"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli/common"
)

type Command struct {
	cmd    string
	params common.Params
}

func New(p common.Params) (*Command, error) {
	if len(p.Args) == 0 {
		return nil, fmt.Errorf("no command specified: %s", help())
	}

	return &Command{
		cmd:    p.Args[0],
		params: p.Next(),
	}, nil
}

func (c *Command) Run() error {
	if command.IsHelp(c.cmd) {
		return fmt.Errorf(help())
	}

	switch c.cmd {
	case "user":
		cc, err := user.New(c.params)
		if err != nil {
			return err
		}

		return cc.Run()
	default:
		return fmt.Errorf("unknown command: %s", c.cmd)
	}
}

func help() string {
	return "available commands: \n" + "user"
}
