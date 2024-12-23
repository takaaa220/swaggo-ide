package handler

import (
	"context"
	"encoding/json"

	"github.com/takaaa220/swaggo-ide/server/internal/handler/protocol"
	"golang.org/x/exp/jsonrpc2"
)

func (h *LSPHandler) HandleDidCloseTextDocument(ctx context.Context, req *jsonrpc2.Request) error {
	var params protocol.DidCloseTextDocumentParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	return h.doDidCloseTextDocument(ctx, &params)
}

func (h *LSPHandler) doDidCloseTextDocument(_ context.Context, p *protocol.DidCloseTextDocumentParams) error {
	h.fileCache.Delete(p.TextDocument.Uri)
	return nil
}
