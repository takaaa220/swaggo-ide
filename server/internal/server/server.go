package server

import (
	"context"
	"log"

	"github.com/takaaa220/swaggo-ide/server/internal/server/transport"
	"golang.org/x/exp/jsonrpc2"
)

func StartServer(ctx context.Context) error {
	handler := NewLSPHandler(LSPHandlerOptions{})

	binder := transport.NewBinder(handler)
	listener := transport.NewStdListener()

	server, err := jsonrpc2.Serve(ctx, listener, binder)
	if err != nil {
		return err
	}

	log.Println("Server started")

	if err := server.Wait(); err != nil {
		return err
	}

	return nil
}
