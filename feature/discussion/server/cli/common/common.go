package common

import (
	"database/sql"
	"io"
)

type Params struct {
	Args []string
	Out  io.Writer
	DB   *sql.DB
}

func (p Params) Next() Params {
	return Params{
		Args: p.Args[1:],
		Out:  p.Out,
		DB:   p.DB,
	}
}
