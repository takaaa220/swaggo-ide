package handler

import (
	"encoding/json"

	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/protocol"
	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/transport"
	"golang.org/x/exp/jsonrpc2"
)

func (h *LSPHandlerOptions) handleCompletion(ctx transport.Context, req *jsonrpc2.Request) (any, error) {
	if h.HandleCompletion == nil {
		return nil, jsonrpc2.ErrNotHandled
	}

	var params protocol.CompletionParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	res, err := h.HandleCompletion(ctx, &params)
	if err != nil {
		return nil, err
	}

	// this is hack because occur error when both of res and err are nil
	if res == nil {
		return []protocol.CompletionItem{}, nil
	}

	return res, nil
}
