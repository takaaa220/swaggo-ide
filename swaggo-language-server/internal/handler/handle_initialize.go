package handler

import (
	"context"
	"encoding/json"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/handler/filecache"
	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/handler/protocol"
	"golang.org/x/exp/jsonrpc2"
)

func (h *LSPHandler) HandleInitialize(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	var params protocol.InitializeParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	return h.doInitialize(ctx, &params)
}

func (h *LSPHandler) doInitialize(_ context.Context, _ *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	h.fileCache = filecache.NewFileCache()

	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			CompletionProvider: &protocol.CompletionOptions{
				TriggerCharacters: []string{" ", "@"},
			},
			HoverProvider: true,
			TextDocumentSync: protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    protocol.TextDocumentSyncKindFull,
				Save: protocol.SaveOptions{
					IncludeText: true,
				},
			},
			CodeLensProvider: &protocol.CodeLensOptions{
				ResolveProvider: true,
			},
		},
	}, nil
}
