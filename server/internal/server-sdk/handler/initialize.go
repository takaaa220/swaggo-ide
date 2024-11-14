package handler

import (
	"encoding/json"

	"github.com/takaaa220/go-swag-ide/server/internal/server-sdk/protocol"
	"github.com/takaaa220/go-swag-ide/server/internal/server-sdk/transport"
	"golang.org/x/exp/jsonrpc2"
)

func (h *LSPHandlerOptions) handleInitialize(ctx transport.Context, req *jsonrpc2.Request) (any, error) {
	if h.HandleInitialize == nil {
		return nil, jsonrpc2.ErrNotHandled
	}

	var params protocol.InitializeParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	return h.HandleInitialize(ctx, &params)
}
