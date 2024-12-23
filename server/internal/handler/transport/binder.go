package transport

import (
	"context"

	"golang.org/x/exp/jsonrpc2"
)

type Binder interface {
	jsonrpc2.Binder
}

type Handler interface {
	jsonrpc2.Handler
	SetConnection(conn *jsonrpc2.Connection)
}

type binder struct {
	handler Handler
}

func NewBinder(handler Handler) *binder {
	return &binder{
		handler: handler,
	}
}

func (b *binder) Bind(ctx context.Context, conn *jsonrpc2.Connection) (jsonrpc2.ConnectionOptions, error) {
	// maybe, this is hack...
	b.handler.SetConnection(conn)

	return jsonrpc2.ConnectionOptions{
		Handler: b.handler,
	}, nil
}
