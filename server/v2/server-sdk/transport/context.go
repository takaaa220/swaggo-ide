package transport

import (
	"context"

	"golang.org/x/exp/jsonrpc2"
)

type Context struct {
	context.Context
	conn *jsonrpc2.Connection
}

func (l *Context) Notify(method string, params any) error {
	return l.conn.Notify(l, method, params)
}
