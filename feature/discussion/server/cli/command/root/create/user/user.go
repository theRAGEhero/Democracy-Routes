package user

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/authenticationhandler"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/cli/common"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/password"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/userhandler"
	"github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/userhandler/model"
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
	fs.StringVar(&c.password, "pass", "", "(optional) User password. If not provided, a random password will be generated.")

	if err := fs.Parse(p.Args); err != nil {
		return nil, fmt.Errorf("parsing flags: %w", err)
	}

	return &c, nil
}

type Response struct {
	ID       string
	Name     string
	Password string
}

func (c *Command) Run() error {
	uh, err := userhandler.New(c.params.DB)
	if err != nil {
		return fmt.Errorf("creating user handler: %w", err)
	}

	user, err := uh.Create(&model.CreateUser{
		Name: c.name,
	})
	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}

	ah, err := authenticationhandler.New(c.params.DB)
	if err != nil {
		return fmt.Errorf("creating auth handler: %w", err)
	}

	pass := c.password
	if pass == "" {
		pass = password.New()
	}

	err = ah.SetPassword(user.ID, pass)
	if err != nil {
		return fmt.Errorf("setting password: %w", err)
	}

	var res Response
	res.ID = user.ID
	res.Name = user.Name
	res.Password = pass

	err = json.NewEncoder(c.params.Out).Encode(res)
	if err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	return nil
}
