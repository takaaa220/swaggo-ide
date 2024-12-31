package internal

import (
	"context"
	"log"

	"github.com/takaaa220/swaggo-ide/server/internal/handler"
	"github.com/takaaa220/swaggo-ide/server/internal/handler/transport"
	"golang.org/x/exp/jsonrpc2"
)

func StartServer(ctx context.Context, debug bool) error {
	opts := handler.LSPHandlerOptions{
		LogLevel: handler.LogWarn,
	}
	if debug {
		opts.LogLevel = handler.LogDebug
	}

	handler := handler.NewLSPHandler(opts)

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
