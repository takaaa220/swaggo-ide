package handler

import (
	"context"
	"encoding/json"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/handler/filecache"
	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/handler/protocol"
	"golang.org/x/exp/jsonrpc2"
)

func (h *LSPHandler) HandleDidChangeTextDocument(ctx context.Context, req *jsonrpc2.Request) error {
	var params protocol.DidChangeTextDocumentParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	return h.doDidChangeTextDocument(ctx, &params)
}

func (h *LSPHandler) doDidChangeTextDocument(ctx context.Context, p *protocol.DidChangeTextDocumentParams) error {
	info, found := h.fileCache.Get(p.TextDocument.Uri)
	if !found {
		return nil
	}

	go func(ctx context.Context) {
		newText := info.Text.Update(p.ContentChanges)
		h.fileCache.Set(p.TextDocument.Uri, filecache.NewFileInfo(p.TextDocument.Version, newText))
	}(ctx)

	return nil
}
