package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/takaaa220/swaggo-ide/server/internal/handler/protocol"
	"github.com/takaaa220/swaggo-ide/server/internal/swag"
	"golang.org/x/exp/jsonrpc2"
)

func (h *LSPHandler) HandleHover(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	var params protocol.HoverParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	res, err := h.doHover(ctx, &params)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *LSPHandler) doHover(_ context.Context, p *protocol.HoverParams) (*protocol.Hover, error) {
	fileInfo, found := h.fileCache.Get(p.TextDocument.Uri)
	if !found {
		return nil, nil
	}

	line, ok := fileInfo.Text.GetLine(int(p.Position.Line))
	if !ok {
		return nil, nil
	}

	attribute, err := swag.GetAttribute(line)
	if err != nil {
		return nil, err
	}

	if attribute == nil {
		return nil, nil
	}

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: fmt.Sprintf("```\n%s\n```\n\n---\n%s", attribute.Title, attribute.Description),
		},
	}, nil
}
