package server

import (
	"context"
	"log"

	"github.com/takaaa220/go-swag-ide/server/internal/server-sdk/handler"
	"github.com/takaaa220/go-swag-ide/server/internal/server-sdk/protocol"
	"github.com/takaaa220/go-swag-ide/server/internal/server-sdk/transport"
	"github.com/takaaa220/go-swag-ide/server/internal/server/swag"
	"github.com/takaaa220/go-swag-ide/server/internal/server/util"
)

var (
	cache *FileCache
)

func StartServer(ctx context.Context) error {
	handler := handler.NewLSPHandler(handler.LSPHandlerOptions{
		HandleInitialize: handleInitialize,
		HandleCompletion: handleCompletion,
		HandleDidOpenTextDocument: func(ctx transport.Context, p *protocol.DidOpenTextDocumentParams) error {
			cache.Set(p.TextDocument.Uri, NewFileInfo(p.TextDocument.Version, NewFileText(p.TextDocument.Text)))

			go func() {
				if err := ctx.Notify("textDocument/publishDiagnostics",
					protocol.PublishDiagnosticsParams{
						Uri:         p.TextDocument.Uri,
						Diagnostics: swag.CheckSyntax(string(p.TextDocument.Uri), p.TextDocument.Text),
					}); err != nil {
					log.Println(err)
				}
			}()

			return nil
		},
		HandleDidChangeTextDocument: func(ctx transport.Context, p *protocol.DidChangeTextDocumentParams) error {
			info, found := cache.Get(p.TextDocument.Uri)
			if !found {
				return nil
			}

			go func() {
				newText := info.Text.Update(p.ContentChanges)

				cache.Set(p.TextDocument.Uri, NewFileInfo(p.TextDocument.Version, newText))
			}()

			return nil
		},
		HandleDidCloseTextDocument: func(ctx transport.Context, p *protocol.DidCloseTextDocumentParams) error {
			cache.Delete(p.TextDocument.Uri)
			return nil
		},
		HandleDidSaveTextDocument: func(ctx transport.Context, p *protocol.DidSaveTextDocumentParams) error {
			cache.Set(p.TextDocument.Uri, NewFileInfo(0, NewFileText(p.Text)))

			log.Println("Saved")

			go func() {
				if err := ctx.Notify("textDocument/publishDiagnostics",
					protocol.PublishDiagnosticsParams{
						Uri:         p.TextDocument.Uri,
						Diagnostics: swag.CheckSyntax(string(p.TextDocument.Uri), p.Text),
					}); err != nil {
					log.Println(err)
				}
			}()

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
			TextDocumentSync: protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    protocol.TextDocumentSyncKindFull,
				Save: protocol.SaveOptions{
					IncludeText: true,
				},
			},
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

	line, ok := fileInfo.Text.GetLine(int(p.Position.Line))
	if !ok {
		return nil, nil
	}
	if !util.IsCommentLine(line) {
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
