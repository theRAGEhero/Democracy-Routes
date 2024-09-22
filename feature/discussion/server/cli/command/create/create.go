package create

import "fmt"

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
	case "user":

		return nil
	default:
		return fmt.Errorf("unknown command: %s", c.cmd)
	}
}
