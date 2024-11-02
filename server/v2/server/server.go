package server

import (
	"context"
	"log"
	"strings"

	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/handler"
	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/protocol"
	"github.com/takaaa220/go-swag-ide/server/v2/server-sdk/transport"
	"github.com/takaaa220/go-swag-ide/server/v2/server/swag"
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
		return nil, err
	}
	if !t {
		return nil, nil
	}

	if p.Context.TriggerCharacter == "@" {
		completionList, err := swag.GetCompletionItems(p.Position)
		if err != nil {
			return nil, err
		}

		return completionList, nil
	}

	return nil, nil
}

func isInFunctionComment(src string, pos protocol.Position) (bool, error) {
	splitSrc := strings.Split(src, "\n")
	if int(pos.Line) >= len(splitSrc) {
		return false, nil
	}

	line := splitSrc[pos.Line]

	return strings.HasPrefix(strings.TrimLeft(line, " \t"), "//"), nil
}
