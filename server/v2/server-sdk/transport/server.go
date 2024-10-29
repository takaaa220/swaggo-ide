package transport

import (
	"context"

	"golang.org/x/exp/jsonrpc2"
)

func NewServer(ctx context.Context, listener jsonrpc2.Listener, binder Binder) (*jsonrpc2.Server, error) {
	server, err := jsonrpc2.Serve(ctx, listener, binder)
	if err != nil {
		return nil, err
	}

	return server, err
}
