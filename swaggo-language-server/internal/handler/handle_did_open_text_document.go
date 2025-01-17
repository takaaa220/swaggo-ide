package handler

import (
	"context"
	"encoding/json"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/handler/filecache"
	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/handler/protocol"
	"golang.org/x/exp/jsonrpc2"
)

func (h *LSPHandler) HandleDidOpenTextDocument(ctx context.Context, req *jsonrpc2.Request) error {
	var params protocol.DidOpenTextDocumentParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	return h.doDidOpenTextDocument(ctx, &params)
}

func (h *LSPHandler) doDidOpenTextDocument(ctx context.Context, p *protocol.DidOpenTextDocumentParams) error {
	h.fileCache.Set(p.TextDocument.Uri, filecache.NewFileInfo(p.TextDocument.Version, filecache.NewFileText(p.TextDocument.Text)))

	h.requestCheckSyntax(p.TextDocument.Uri, p.TextDocument.Text)

	return nil
}
