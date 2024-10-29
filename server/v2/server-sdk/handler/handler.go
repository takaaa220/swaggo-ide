package handler

import (
	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/transport"
	"golang.org/x/exp/jsonrpc2"
)

type LSPHandlerOptions struct {
	HandleInitialize InitializeFunc
}

func (h *LSPHandlerOptions) Handle(ctx transport.Context, req *jsonrpc2.Request) (any, error) {
	switch req.Method {
	case "initialize":
		return h.handleInitialize(ctx, req)
	default:
		return nil, jsonrpc2.ErrNotHandled
	}
}

func NewLSPHandler(opts LSPHandlerOptions) *transport.BaseHandler {
	return transport.NewBaseHandler(opts.Handle)
}
