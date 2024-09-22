package common

import "io"

type Params struct {
	Args []string
	Out  io.Writer
}

func (p Params) Next() Params {
	return Params{
		Args: p.Args[1:],
		Out:  p.Out,
	}
}
