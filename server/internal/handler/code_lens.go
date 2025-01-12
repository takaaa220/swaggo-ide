package handler

import (
	"context"
	"encoding/json"

	"github.com/takaaa220/swaggo-ide/server/internal/handler/protocol"
	"github.com/takaaa220/swaggo-ide/server/internal/handler/swag"
	"golang.org/x/exp/jsonrpc2"
)

func (h *LSPHandler) HandleCodeLens(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	var params protocol.CodeLensParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	res, err := h.doCodeLens(ctx, &params)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *LSPHandler) doCodeLens(_ context.Context, p *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	fileInfo, found := h.fileCache.Get(p.TextDocument.Uri)
	if !found {
		return nil, nil
	}

	swagCommentsRanges := swag.FindSwagComments(fileInfo.Text.String())
	codeLens := make([]protocol.CodeLens, len(swagCommentsRanges))
	for i, swagCommentsRange := range swagCommentsRanges {
		codeLens[i] = protocol.CodeLens{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(swagCommentsRange.Start),
					Character: uint32(0),
				},
				End: protocol.Position{
					Line:      uint32(swagCommentsRange.Start), // should only span a single line, so start and end are the same
					Character: uint32(0),
				},
			},
			Command: protocol.Command{
				Title:   "swag fmt (this file)",
				Command: "swaggo-language-server-client.format",
				Arguments: []any{
					p.TextDocument.Uri,
				},
			},
		}
	}

	return codeLens, nil
}
