package handler

import (
	"context"
	"encoding/json"

	"github.com/takaaa220/swaggo-ide/server/internal/handler/protocol"
	"github.com/takaaa220/swaggo-ide/server/internal/swag"
	"golang.org/x/exp/jsonrpc2"
)

func (h *LSPHandler) HandleCompletion(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	var params protocol.CompletionParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	res, err := h.doCompletion(ctx, &params)
	if err != nil {
		return nil, err
	}

	// this is hack because occur error when both of res and err are nil
	if res == nil {
		return []protocol.CompletionItem{}, nil
	}

	return res, nil
}

func (h *LSPHandler) doCompletion(_ context.Context, p *protocol.CompletionParams) (protocol.CompletionResult, error) {
	fileInfo, found := h.fileCache.Get(p.TextDocument.Uri)
	if !found {
		return nil, nil
	}

	line, ok := fileInfo.Text.GetLine(int(p.Position.Line))
	if !ok {
		return nil, nil
	}

	switch p.Context.TriggerCharacter {
	case "@":
		return swag.GetTagCompletionItems(line, p.Position)
	case " ":
		return swag.GetTagArgCompletionItems(line, p.Position)
	default:
		return nil, nil
	}
}
