package handler

import (
	"context"

	"github.com/takaaa220/swaggo-ide/server/internal/handler/filecache"
	"golang.org/x/exp/jsonrpc2"
)

func NewLSPHandler(opts LSPHandlerOptions) *LSPHandler {
	return &LSPHandler{
		opts: opts,
	}
}

type LSPHandler struct {
	opts      LSPHandlerOptions
	conn      *jsonrpc2.Connection
	fileCache *filecache.FileCache
}

type LSPHandlerOptions struct{}

func (h *LSPHandler) Handle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	switch req.Method {
	case "initialize":
		return h.HandleInitialize(ctx, req)
	case "textDocument/didOpen":
		err := h.HandleDidOpenTextDocument(ctx, req)
		if err != nil {
			return nil, err
		}
		return nil, jsonrpc2.ErrAsyncResponse
	case "textDocument/didChange":
		err := h.HandleDidChangeTextDocument(ctx, req)
		if err != nil {
			return nil, err
		}
		return nil, jsonrpc2.ErrAsyncResponse
	case "textDocument/didClose":
		err := h.HandleDidCloseTextDocument(ctx, req)
		if err != nil {
			return nil, err
		}
		return nil, jsonrpc2.ErrAsyncResponse
	case "textDocument/didSave":
		err := h.HandleDidSaveTextDocument(ctx, req)
		if err != nil {
			return nil, err
		}
		return nil, jsonrpc2.ErrAsyncResponse
	case "textDocument/completion":
		return h.HandleCompletion(ctx, req)
	default:
		return nil, jsonrpc2.ErrNotHandled
	}
}

func (h *LSPHandler) SetConnection(conn *jsonrpc2.Connection) {
	h.conn = conn
}

func (h *LSPHandler) Notify(ctx context.Context, method string, params any) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return h.conn.Notify(ctx, method, params)
	}
}
