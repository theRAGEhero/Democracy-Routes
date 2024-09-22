package user

import (
	"flag"
	"fmt"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli/common"
)

type Command struct {
	params common.Params

	name     string
	password string
}

func New(p common.Params) (*Command, error) {
	fs := flag.NewFlagSet("user", flag.ContinueOnError)

	var c Command
	c.params = p

	fs.StringVar(&c.name, "name", "", "User name.")
	fs.StringVar(&c.password, "password", "", "(optional) User password. If not provided, a random password will be generated.")

	if err := fs.Parse(p.Args); err != nil {
		return nil, fmt.Errorf("parsing flags: %w", err)
	}

	return &c, nil
}

func (c *Command) Run() error {
	_, err := c.params.Out.Write([]byte("created"))
	if err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	return nil
}
