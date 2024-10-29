package handler

import (
	"encoding/json"

	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/protocol"
	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/transport"
	"golang.org/x/exp/jsonrpc2"
)

type InitializeFunc func(transport.Context, *protocol.InitializeParams) (*protocol.InitializeResult, error)

func (h *LSPHandlerOptions) handleInitialize(ctx transport.Context, req *jsonrpc2.Request) (any, error) {
	if h.HandleInitialize == nil {
		return nil, jsonrpc2.ErrNotHandled
	}

	var params protocol.InitializeParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	_, err := h.HandleInitialize(ctx, &params)
	if err != nil {
		return nil, err
	}

	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{},
	}, nil
}
