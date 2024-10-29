package server

import (
	"context"

	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/handler"
	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/protocol"
	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/transport"
)

func StartServer(ctx context.Context) error {
	handler := handler.NewLSPHandler(handler.LSPHandlerOptions{
		HandleInitialize: i,
	})

	binder := transport.NewStdioBinder(handler)
	listener := transport.NewStdListener()

	server, err := transport.NewServer(ctx, listener, binder)
	if err != nil {
		return err
	}

	if err := server.Wait(); err != nil {
		return err
	}

	return nil
}

func i(ctx transport.Context, p *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	ctx.Notify("window/logMessage", protocol.LogMessageParams{
		Type:    protocol.MessageTypeInfo,
		Message: "Initializing Language Server",
	})

	return nil, nil
}
