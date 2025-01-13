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
		return Null{}, nil
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
		candidates, err := swag.GetTagCompletionItems(line)
		if err != nil {
			return nil, err
		}
		return convertCandidates(
			candidates,
			protocol.Range{
				Start: protocol.Position{Line: p.Position.Line, Character: p.Position.Character - 1},
				End:   p.Position,
			},
		), nil
	case " ":
		candidates, err := swag.GetTagArgCompletionItems(line, convertPosition(p.Position))
		if err != nil {
			return nil, err
		}
		return convertCandidates(
			candidates,
			protocol.Range{
				Start: p.Position,
				End:   p.Position,
			},
		), nil
	default:
		return nil, nil
	}
}

func convertPosition(position protocol.Position) swag.Position {
	return swag.Position{
		Line:      position.Line,
		Character: position.Character,
	}
}

func convertCandidates(candidates []swag.CompletionCandidate, textEditRange protocol.Range) *protocol.CompletionList {
	if len(candidates) == 0 {
		return nil
	}

	completionItems := make([]protocol.CompletionItem, len(candidates))
	for i, candidate := range candidates {
		completionItems[i] = protocol.CompletionItem{
			Label: candidate.Label,
			Kind:  protocol.CompletionItemKindKeyword,
			TextEdit: protocol.TextEdit{
				Range:   textEditRange,
				NewText: candidate.NewText,
			},
		}
	}

	return &protocol.CompletionList{
		Items:        completionItems,
		IsIncomplete: false,
	}
}
