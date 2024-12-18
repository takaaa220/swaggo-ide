package handler

import (
	"fmt"

	"github.com/takaaa220/swaggo-ide/server/internal/server-sdk/protocol"
	"github.com/takaaa220/swaggo-ide/server/internal/server-sdk/transport"
	"golang.org/x/exp/jsonrpc2"
)

type LSPHandlerOptions struct {
	HandleInitialize protocol.InitializeFunc
	HandleCompletion protocol.CompletionFunc

	HandleDidOpenTextDocument   protocol.DidOpenTextDocumentFunc
	HandleDidChangeTextDocument protocol.DidChangeTextDocumentFunc
	HandleDidCloseTextDocument  protocol.TextDocumentDidCloseFunc
	HandleDidSaveTextDocument   protocol.TextDocumentDidSaveFunc
}

func (h *LSPHandlerOptions) Handle(ctx transport.Context, req *jsonrpc2.Request) (any, error) {
	switch req.Method {
	case "initialize":
		return h.handleInitialize(ctx, req)
	case "textDocument/didOpen":
		err := h.handleDidOpenTextDocument(ctx, req)
		if err != nil {
			return nil, err
		}
		return nil, jsonrpc2.ErrAsyncResponse
	case "textDocument/didChange":
		err := h.handleDidChangeTextDocument(ctx, req)
		if err != nil {
			return nil, err
		}
		return nil, jsonrpc2.ErrAsyncResponse
	case "textDocument/didClose":
		err := h.handleDidCloseTextDocument(ctx, req)
		if err != nil {
			return nil, err
		}
		return nil, jsonrpc2.ErrAsyncResponse
	case "textDocument/didSave":
		err := h.handleDidSaveTextDocument(ctx, req)
		if err != nil {
			return nil, err
		}
		return nil, jsonrpc2.ErrAsyncResponse
	case "textDocument/completion":
		return h.handleCompletion(ctx, req)
	default:
		protocol.NewLogger(protocol.MessageTypeLog).Error(ctx, fmt.Sprintf("unknown method: %s", req.Method))
		return nil, jsonrpc2.ErrNotHandled
	}
}

func NewLSPHandler(opts LSPHandlerOptions) *transport.BaseHandler {
	return transport.NewBaseHandler(opts.Handle)
}
