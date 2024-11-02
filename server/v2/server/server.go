package server

import (
	"context"
	"log"

	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/handler"
	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/protocol"
	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/transport"
)

var (
	cache *FileCache
)

func StartServer(ctx context.Context) error {
	handler := handler.NewLSPHandler(handler.LSPHandlerOptions{
		HandleInitialize: handleInitialize,
		HandleCompletion: handleCompletion,
		HandleDidOpenTextDocument: func(ctx transport.Context, p *protocol.DidOpenTextDocumentParams) error {
			log.Println("handleDidOpenTextDocument", p.TextDocument.Uri)
			cache.Set(p.TextDocument.Uri, NewFileInfo(p.TextDocument.Version, NewFileText(p.TextDocument.Text)))
			return nil
		},
		HandleDidChangeTextDocument: func(ctx transport.Context, p *protocol.DidChangeTextDocumentParams) error {
			log.Println("handleDidChangeTextDocument", p.TextDocument.Uri)
			info, found := cache.Get(p.TextDocument.Uri)
			if !found {
				return nil
			}

			newText := info.Text.Update(p.ContentChanges)

			cache.Set(p.TextDocument.Uri, NewFileInfo(p.TextDocument.Version, newText))

			return nil
		},
		HandleDidCloseTextDocument: func(ctx transport.Context, p *protocol.DidCloseTextDocumentParams) error {
			log.Println("handleDidCloseTextDocument", p.TextDocument.Uri)
			cache.Delete(p.TextDocument.Uri)
			return nil
		},
		HandleDidSaveTextDocument: func(ctx transport.Context, p *protocol.DidSaveTextDocumentParams) error {
			log.Println("handleDidSaveTextDocument", p.TextDocument.Uri)
			cache.Delete(p.TextDocument.Uri)
			cache.Set(p.TextDocument.Uri, NewFileInfo(0, NewFileText(p.Text)))
			return nil
		},
	})

	binder := transport.NewStdioBinder(handler)
	listener := transport.NewStdListener()

	server, err := transport.NewServer(ctx, listener, binder)
	if err != nil {
		return err
	}

	log.Println("Server started")

	if err := server.Wait(); err != nil {
		return err
	}

	return nil
}

func handleInitialize(ctx transport.Context, p *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	if err := protocol.NewLogger(protocol.MessageTypeLog).Info(ctx, "initialize"); err != nil {
		return nil, err
	}

	cache = NewFileCache()

	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			CompletionProvider: &protocol.CompletionOptions{
				TriggerCharacters: []string{" ", "@"},
			},
			TextDocumentSync: protocol.TextDocumentSyncKindFull,
		},
	}, nil
}

func handleCompletion(ctx transport.Context, p *protocol.CompletionParams) (protocol.CompletionResult, error) {
	if err := protocol.NewLogger(protocol.MessageTypeLog).Info(ctx, "completion"); err != nil {
		return nil, err
	}

	fileInfo, found := cache.Get(p.TextDocument.Uri)
	if !found {
		return nil, nil
	}

	t, err := isInFunctionComment(fileInfo.Text.String(), p.Position)
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
