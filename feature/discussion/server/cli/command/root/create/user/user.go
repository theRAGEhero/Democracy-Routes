package user

import (
	"flag"
	"fmt"
)

type Command struct {
	name     string
	password string
}

func New(args []string) (*Command, error) {
	fs := flag.NewFlagSet("user", flag.ContinueOnError)

	var c Command

	fs.StringVar(&c.name, "name", "", "User name.")
	fs.StringVar(&c.password, "password", "", "(optional) User password. If not provided, a random password will be generated.")

	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("parsing flags: %w", err)
	}

	return &Command{}, nil
}

func (c *Command) Run() error {
	return nil
}
