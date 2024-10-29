package server

import (
	"context"
	"go/token"
	"log"
	"net/url"

	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/handler"
	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/protocol"
	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/transport"
)

func StartServer(ctx context.Context) error {
	handler := handler.NewLSPHandler(handler.LSPHandlerOptions{
		HandleInitialize: handleInitialize,
		HandleCompletion: handleCompletion,
	})

	binder := transport.NewStdioBinder(handler)
	listener := transport.NewStdListener()

	server, err := transport.NewServer(ctx, listener, binder)
	if err != nil {
		return err
	}

	if err := server.Wait(); err != nil {
		return err
	}

	return nil
}

func handleInitialize(ctx transport.Context, p *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	if err := protocol.NewLogger(protocol.MessageTypeLog).Info(ctx, "initialize"); err != nil {
		return nil, err
	}

	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			CompletionProvider: &protocol.CompletionOptions{
				TriggerCharacters: []string{" ", "@"},
			},
		},
	}, nil
}

func handleCompletion(ctx transport.Context, p *protocol.CompletionParams) (protocol.CompletionResult, error) {
	if err := protocol.NewLogger(protocol.MessageTypeLog).Info(ctx, "completion"); err != nil {
		return nil, err
	}

	parsedUri, err := url.Parse(string(p.TextDocument.Uri))
	if err != nil {
		log.Println("url.Parse", err)
		return nil, err
	}

	log.Println("parsedUri.Path", parsedUri.Path)

	t, err := isInFunctionComment(token.Position{
		Filename: parsedUri.Path,
		Line:     int(p.Position.Line + 1),
		Column:   int(p.Position.Character + 1),
		Offset:   -1,
	})
	if err != nil {
		log.Println("isInFunctionComment", err)
		return nil, err
	}

	log.Println("isInFunctionComment", t)

	if !t {
		return nil, nil
	}
	if p.Context.TriggerCharacter == "@" {
		completionList, err := getCompletionItems(p.Position)
		if err != nil {
			return nil, err
		}

		return completionList, nil
	}

	return []protocol.CompletionItem{
		{
			Label: "Hello1",
			Kind:  protocol.CompletionItemKindKeyword,
		},
		{
			Label: "Hello2",
			Kind:  protocol.CompletionItemKindKeyword,
		},
	}, nil
}

type swagTag struct {
	label string
	args  []string
}

var swagTags = []swagTag{
	{
		label: "@Summary",
		args:  []string{"SUMMARY"},
	},
	{
		label: "@Description",
		args:  []string{"DESCRIPTION"},
	},
	{
		label: "@Tags",
		args:  []string{"TAG1,TAG2"},
	},
	{
		label: "@Accept",
		args:  []string{"MIME_TYPE"},
	},
	{
		label: "@Produce",
		args:  []string{"MIME_TYPE"},
	},
	{
		label: "@Param",
		args:  []string{"PARAM_NAME", "PARAM_TYPE", "DATA_TYPE", "REQUIRED(bool)", "DESCRIPTION", "ATTRIBUTE(optional)"},
	},
	{
		label: "@Success",
		args:  []string{"STATUS_CODE", "{DATA_TYPE}", "DESCRIPTION"},
	},
	{
		label: "@Failure",
		args:  []string{"STATUS_CODE", "{DATA_TYPE}", "DESCRIPTION"},
	},
	{
		label: "@Router",
		args:  []string{"PATH", "[METHOD]"},
	},
	{
		label: "@Security",
		args:  []string{},
	},
	{
		label: "@ID",
		args:  []string{"ID"},
	},
	{
		label: "@Header",
		args:  []string{"STATUS_CODE", "{PARAM_TYPE}", "DATA_TYPE", "COMMENT"},
	},
}

func getCompletionItems(position protocol.Position) (*protocol.CompletionList, error) {
	kind := protocol.CompletionItemKindKeyword
	completionItems := make([]protocol.CompletionItem, len(swagTags))
	for i, tag := range swagTags {
		completionText := tag.label
		for _, arg := range tag.args {
			completionText += "  " + arg
		}

		completionItems[i] = protocol.CompletionItem{
			Label: tag.label,
			Kind:  kind,
			TextEdit: &protocol.TextEdit{
				NewText: completionText,
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      position.Line,
						Character: position.Character - 1,
					},
					End: position,
				},
			},
		}

	}

	return &protocol.CompletionList{
		IsIncomplete: false,
		Items:        completionItems,
	}, nil
}
