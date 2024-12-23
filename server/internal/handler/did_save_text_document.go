package handler

import (
	"context"
	"encoding/json"

	"github.com/takaaa220/swaggo-ide/server/internal/handler/filecache"
	"github.com/takaaa220/swaggo-ide/server/internal/handler/protocol"
	"golang.org/x/exp/jsonrpc2"
)

func (h *LSPHandler) HandleDidSaveTextDocument(ctx context.Context, req *jsonrpc2.Request) error {
	var params protocol.DidSaveTextDocumentParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	return h.doDidSaveTextDocument(ctx, &params)
}

func (h *LSPHandler) doDidSaveTextDocument(_ context.Context, p *protocol.DidSaveTextDocumentParams) error {
	h.fileCache.Set(p.TextDocument.Uri, filecache.NewFileInfo(0, filecache.NewFileText(p.Text)))

	h.logger.Debugf("Saved: %s", p.TextDocument.Uri)

	h.requestCheckSyntax(p.TextDocument.Uri, p.Text)

	return nil
}
