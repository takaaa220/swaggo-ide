package handler

import (
	"context"
	"encoding/json"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/handler/protocol"
	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/swag"
	"golang.org/x/exp/jsonrpc2"
)

func (h *LSPHandler) HandleSignatureHelp(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	var params protocol.SignatureHelpParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, protocol.NewResponseError(protocol.CodeInvalidParams, err.Error(), nil)
	}

	return h.doSignatureHelp(ctx, &params)
}

func (h *LSPHandler) doSignatureHelp(_ context.Context, params *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) {
	fileInfo, found := h.fileCache.Get(params.TextDocument.Uri)
	if !found {
		return nil, nil
	}

	line, ok := fileInfo.Text.GetLine(int(params.Position.Line))
	if !ok {
		return nil, nil
	}

	switch params.Context.TriggerCharacter {
	case " ":
		res, err := swag.GetTagArg(line, convertPosition(params.Position))
		if err != nil {
			return nil, err
		}
		if res == nil {
			return nil, nil
		}

		var (
			activeSignature uint = 0
			activeParameter uint = uint(res.ActiveArgIndex)
		)

		parameters := make([]protocol.ParameterInformation, len(res.Tag.Args))
		for i, arg := range res.Tag.Args {
			parameters[i] = protocol.ParameterInformation{
				Label: arg.Label(),
			}
		}

		return &protocol.SignatureHelp{
			Signatures: []protocol.SignatureInformation{
				{
					Label:      res.Tag.String(),
					Parameters: parameters,
					Documentation: &protocol.MarkupContent{
						Kind:  protocol.MarkupKindMarkdown,
						Value: res.Tag.Description,
					},
				},
			},
			ActiveSignature: &activeSignature,
			ActiveParameter: &activeParameter,
		}, nil
	default:
		return nil, nil
	}
}
