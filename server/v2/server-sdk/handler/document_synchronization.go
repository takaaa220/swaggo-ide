package handler

import (
	"encoding/json"

	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/protocol"
	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/transport"
	"golang.org/x/exp/jsonrpc2"
)

func (h *LSPHandlerOptions) handleDidOpenTextDocument(ctx transport.Context, req *jsonrpc2.Request) error {
	if h.HandleDidOpenTextDocument == nil {
		return jsonrpc2.ErrNotHandled
	}

	var params protocol.DidOpenTextDocumentParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	return h.HandleDidOpenTextDocument(ctx, &params)
}

func (h *LSPHandlerOptions) handleDidCloseTextDocument(ctx transport.Context, req *jsonrpc2.Request) error {
	if h.HandleDidCloseTextDocument == nil {
		return jsonrpc2.ErrNotHandled
	}

	var params protocol.DidCloseTextDocumentParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	return h.HandleDidCloseTextDocument(ctx, &params)
}

func (h *LSPHandlerOptions) handleDidChangeTextDocument(ctx transport.Context, req *jsonrpc2.Request) error {
	if h.HandleDidChangeTextDocument == nil {
		return jsonrpc2.ErrNotHandled
	}

	var params protocol.DidChangeTextDocumentParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	return h.HandleDidChangeTextDocument(ctx, &params)
}

func (h *LSPHandlerOptions) handleDidSaveTextDocument(ctx transport.Context, req *jsonrpc2.Request) error {
	if h.HandleDidSaveTextDocument == nil {
		return jsonrpc2.ErrNotHandled
	}

	var params protocol.DidSaveTextDocumentParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	return h.HandleDidSaveTextDocument(ctx, &params)
}
