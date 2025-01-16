package handler

import (
	"context"
	"encoding/json"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/handler/protocol"
	"golang.org/x/exp/jsonrpc2"
)

func (h *LSPHandler) HandleCancelRequest(ctx context.Context, req *jsonrpc2.Request) error {
	var params protocol.CancelParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	return h.doCancel(ctx, &params)
}

func (h *LSPHandler) doCancel(_ context.Context, p *protocol.CancelParams) error {
	h.conn.Cancel(p.ID)
	return nil
}
