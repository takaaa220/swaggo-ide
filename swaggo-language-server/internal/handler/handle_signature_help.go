package handler

import (
	"context"
	"encoding/json"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/handler/protocol"
	"golang.org/x/exp/jsonrpc2"
)

func (h *LSPHandler) HandleSignatureHelp(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	var params protocol.SignatureHelpParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	return h.doSignatureHelp(ctx, &params)
}

func (h *LSPHandler) doSignatureHelp(_ context.Context, params *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) {
	return nil, nil
}
