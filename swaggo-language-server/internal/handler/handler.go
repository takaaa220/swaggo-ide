package handler

import (
	"context"
	"io"
	"time"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/handler/filecache"
	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/handler/protocol"
	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/swag"
	"golang.org/x/exp/jsonrpc2"
)

func NewLSPHandler(ctx context.Context, shutdownChan chan struct{}, opts LSPHandlerOptions) *LSPHandler {
	h := &LSPHandler{
		checkSyntaxDebounce: opts.CheckSyntaxDebounce,
		checkSyntaxReq:      make(chan CheckSyntaxRequest),
		logger:              NewLogger(opts.LogWriter, opts.LogLevel),
		shutdownChan:        shutdownChan,
	}

	go h.checkSyntax(ctx)
	return h
}

type LSPHandler struct {
	logger              *logger
	conn                *jsonrpc2.Connection
	fileCache           *filecache.FileCache
	checkSyntaxReq      chan CheckSyntaxRequest
	checkSyntaxDebounce time.Duration
	checkSyntaxTimer    *time.Timer
	shutdownChan        chan struct{}
}

type LSPHandlerOptions struct {
	CheckSyntaxDebounce time.Duration
	LogLevel            LogLevel
	LogWriter           io.Writer
}

func (h *LSPHandler) Handle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	ctx, cancel := h.withHandleTimeout(ctx, req, 5*time.Second)
	defer cancel()

	switch req.Method {
	case "initialize":
		return h.HandleInitialize(ctx, req)
	case "initialized":
		return Null{}, nil
	case "$/cancelRequest":
		return Null{}, h.HandleCancelRequest(ctx, req)
	case "shutdown":
		return Null{}, h.HandleShutdown(ctx, req)
	case "exit":
		h.logger.Debugf("exit received, %v", req.ID)
		return Null{}, nil
	case "textDocument/didOpen":
		err := h.HandleDidOpenTextDocument(ctx, req)
		if err != nil {
			return nil, err
		}
		return Null{}, nil
	case "textDocument/didChange":
		err := h.HandleDidChangeTextDocument(ctx, req)
		if err != nil {
			return nil, err
		}
		return Null{}, nil
	case "textDocument/didClose":
		err := h.HandleDidCloseTextDocument(ctx, req)
		if err != nil {
			return nil, err
		}
		return Null{}, nil
	case "textDocument/didSave":
		err := h.HandleDidSaveTextDocument(ctx, req)
		if err != nil {
			return nil, err
		}
		return Null{}, nil
	case "textDocument/completion":
		return h.HandleCompletion(ctx, req)
	case "textDocument/codeLens":
		return h.HandleCodeLens(ctx, req)
	case "textDocument/hover":
		return h.HandleHover(ctx, req)
	case "workspace/didChangeWatchedFiles":
		// TODO: implement
		return Null{}, nil
	case "$/setTrace":
		// TODO: implement
		return Null{}, nil
	default:
		h.logger.Debugf("method %s not supported", req.Method)
		return nil, jsonrpc2.ErrNotHandled
	}
}

func (h *LSPHandler) withHandleTimeout(ctx context.Context, req *jsonrpc2.Request, timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func(id jsonrpc2.ID) {
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			h.logger.Debugf("context done: %v", ctx.Err())
			h.conn.Cancel(id)
		}
	}(req.ID)

	return ctx, cancel
}

func (h *LSPHandler) SetConnection(conn *jsonrpc2.Connection) {
	h.conn = conn
}

func (h *LSPHandler) CloseConnection() error {
	if h.conn == nil {
		return nil
	}
	return h.conn.Close()
}

func (h *LSPHandler) Notify(ctx context.Context, method string, params any) error {
	if h.conn == nil {
		return nil
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}

	return h.conn.Notify(ctx, method, params)
}

type CheckSyntaxRequest struct {
	uri  protocol.DocumentUri
	text string
}

func (h *LSPHandler) requestCheckSyntax(uri protocol.DocumentUri, text string) {
	if h.checkSyntaxTimer != nil {
		h.checkSyntaxTimer.Reset(h.checkSyntaxDebounce)
	}

	h.checkSyntaxTimer = time.AfterFunc(h.checkSyntaxDebounce, func() {
		h.checkSyntaxTimer = nil
		h.checkSyntaxReq <- CheckSyntaxRequest{
			uri:  uri,
			text: text,
		}
	})
}

func (h *LSPHandler) checkSyntax(ctx context.Context) {
	runningCheckSyntax := make(map[protocol.DocumentUri]context.CancelFunc)

	for {
		select {
		case <-ctx.Done():
			// FIXME: refactor (shouldn't stop timer here)
			for uri, cancel := range runningCheckSyntax {
				cancel()
				delete(runningCheckSyntax, uri)
			}
			if h.checkSyntaxTimer != nil {
				h.checkSyntaxTimer.Stop()
			}
			close(h.checkSyntaxReq)

			return
		case req, ok := <-h.checkSyntaxReq:
			if !ok {
				break
			}

			if cancel, ok := runningCheckSyntax[req.uri]; ok {
				cancel()
			}

			ctx, cancel := context.WithCancel(context.Background())
			runningCheckSyntax[req.uri] = cancel

			go func(uri protocol.DocumentUri, text string) {
				defer cancel()

				syntaxErrors := swag.CheckSyntax(string(uri), text)

				diagnostics := make([]protocol.Diagnostics, len(syntaxErrors))
				for i, syntaxError := range syntaxErrors {
					diagnostics[i] = protocol.Diagnostics{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      syntaxError.Range.Start.Line,
								Character: syntaxError.Range.Start.Character,
							},
							End: protocol.Position{
								Line:      syntaxError.Range.End.Line,
								Character: syntaxError.Range.End.Character,
							},
						},
						Severity: protocol.DiagnosticsSeverityError,
						Source:   "swag",
						Message:  syntaxError.Message,
					}
				}

				if err := h.Notify(ctx, "textDocument/publishDiagnostics",
					protocol.PublishDiagnosticsParams{
						Uri:         uri,
						Diagnostics: diagnostics,
					}); err != nil {
					h.logger.Error(err)
				}
			}(req.uri, req.text)
		}
	}
}

// this is hack (error occurs when Handle returns (nil, nil))
type Null struct{}

func (n *Null) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}
