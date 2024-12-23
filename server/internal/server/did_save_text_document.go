package server

import (
	"context"
	"encoding/json"
	"log"

	"github.com/takaaa220/swaggo-ide/server/internal/server/filecache"
	"github.com/takaaa220/swaggo-ide/server/internal/server/protocol"
	"github.com/takaaa220/swaggo-ide/server/internal/server/swag"
	"golang.org/x/exp/jsonrpc2"
)

func (h *LSPHandler) HandleDidSaveTextDocument(ctx context.Context, req *jsonrpc2.Request) error {
	var params protocol.DidSaveTextDocumentParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	return h.doDidSaveTextDocument(ctx, &params)
}

func (h *LSPHandler) doDidSaveTextDocument(ctx context.Context, p *protocol.DidSaveTextDocumentParams) error {
	h.fileCache.Set(p.TextDocument.Uri, filecache.NewFileInfo(0, filecache.NewFileText(p.Text)))

	log.Println("Saved")

	go func() {
		if err := h.Notify(ctx, "textDocument/publishDiagnostics",
			protocol.PublishDiagnosticsParams{
				Uri:         p.TextDocument.Uri,
				Diagnostics: swag.CheckSyntax(string(p.TextDocument.Uri), p.Text),
			}); err != nil {
			log.Println(err)
		}
	}()

	return nil
}
