package transport

import (
	"context"

	"golang.org/x/exp/jsonrpc2"
)

type HandleFunc func(Context, *jsonrpc2.Request) (any, error)

type BaseHandler struct {
	handleFunc HandleFunc
	conn       *jsonrpc2.Connection
}

func NewBaseHandler(handleFunc HandleFunc) *BaseHandler {
	return &BaseHandler{
		handleFunc: handleFunc,
	}
}

func (h *BaseHandler) Handle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	return h.handleFunc(Context{ctx, h.conn}, req)
}

func (h *BaseHandler) setConnection(conn *jsonrpc2.Connection) {
	h.conn = conn
}

type Binder interface {
	jsonrpc2.Binder
}

type stdioBinder struct {
	handler *BaseHandler
}

var _ Binder = (*stdioBinder)(nil)

func NewStdioBinder(handler *BaseHandler) jsonrpc2.Binder {
	return &stdioBinder{
		handler: handler,
	}
}

func (b *stdioBinder) Bind(ctx context.Context, conn *jsonrpc2.Connection) (jsonrpc2.ConnectionOptions, error) {
	b.handler.setConnection(conn)

	return jsonrpc2.ConnectionOptions{
		Handler: b.handler,
	}, nil
}
