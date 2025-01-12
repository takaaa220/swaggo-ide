package handler

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/takaaa220/swaggo-ide/server/internal/handler/filecache"
	"github.com/takaaa220/swaggo-ide/server/internal/handler/protocol"
	"github.com/takaaa220/swaggo-ide/server/internal/handler/swag"
	"golang.org/x/exp/jsonrpc2"
)

func NewLSPHandler(opts LSPHandlerOptions) *LSPHandler {
	var logWriter io.Writer = os.Stderr
	if opts.LogWriter != nil {
		logWriter = opts.LogWriter
	}

	checkSyntaxDebounce := 100 * time.Millisecond
	if opts.CheckSyntaxDebounce > 0 {
		checkSyntaxDebounce = opts.CheckSyntaxDebounce
	}

	h := &LSPHandler{
		checkSyntaxDebounce: checkSyntaxDebounce,
		checkSyntaxReq:      make(chan CheckSyntaxRequest),
		logger:              NewLogger(logWriter, opts.LogLevel),
	}

	go h.checkSyntax()
	return h
}

type LSPHandler struct {
	logger              *logger
	conn                *jsonrpc2.Connection
	fileCache           *filecache.FileCache
	checkSyntaxReq      chan CheckSyntaxRequest
	checkSyntaxDebounce time.Duration
	checkSyntaxTimer    *time.Timer
}

type LSPHandlerOptions struct {
	CheckSyntaxDebounce time.Duration
	LogLevel            logLevel
	LogWriter           io.Writer
}

func (h *LSPHandler) Handle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	switch req.Method {
	case "initialize":
		return h.HandleInitialize(ctx, req)
	case "initialized":
		return nil, nil
	case "shutdown":
		return h.HandleShutdown(ctx, req)
	case "textDocument/didOpen":
		err := h.HandleDidOpenTextDocument(ctx, req)
		if err != nil {
			return nil, err
		}
		return nil, jsonrpc2.ErrAsyncResponse
	case "textDocument/didChange":
		err := h.HandleDidChangeTextDocument(ctx, req)
		if err != nil {
			return nil, err
		}
		return nil, jsonrpc2.ErrAsyncResponse
	case "textDocument/didClose":
		err := h.HandleDidCloseTextDocument(ctx, req)
		if err != nil {
			return nil, err
		}
		return nil, jsonrpc2.ErrAsyncResponse
	case "textDocument/didSave":
		err := h.HandleDidSaveTextDocument(ctx, req)
		if err != nil {
			return nil, err
		}
		return nil, jsonrpc2.ErrAsyncResponse
	case "textDocument/completion":
		return h.HandleCompletion(ctx, req)
	case "textDocument/codeLens":
		return h.HandleCodeLens(ctx, req)
	case "textDocument/hover":
		return h.HandleHover(ctx, req)
	default:
		h.logger.Debugf("method %s not supported", req.Method)
		return nil, jsonrpc2.ErrNotHandled
	}
}

func (h *LSPHandler) SetConnection(conn *jsonrpc2.Connection) {
	// h.logger.Println("SetConnection called") // bug: this is called multiple times

	h.conn = conn
}

func (h *LSPHandler) CloseConnection() error {
	h.logger.Debugf("CloseConnection")

	if h.conn == nil {
		return nil
	}

	return h.conn.Close()
}

func (h *LSPHandler) Notify(ctx context.Context, method string, params any) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return h.conn.Notify(ctx, method, params)
	}
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

func (h *LSPHandler) checkSyntax() {
	runningCheckSyntax := make(map[protocol.DocumentUri]context.CancelFunc)

	for {
		select {
		case req, ok := <-h.checkSyntaxReq:
			if !ok {
				break
			}

			if cancel, ok := runningCheckSyntax[req.uri]; ok {
				cancel()
			}

			ctx, cancel := context.WithCancel(context.Background())
			runningCheckSyntax[req.uri] = cancel

			go func(ctx context.Context, uri protocol.DocumentUri, text string) {
				if err := h.Notify(ctx, "textDocument/publishDiagnostics",
					protocol.PublishDiagnosticsParams{
						Uri:         uri,
						Diagnostics: swag.CheckSyntax(string(uri), text),
					}); err != nil {
					h.logger.Error(err)
				}
			}(ctx, req.uri, req.text)
		}
	}
}

func (h *LSPHandler) logMessage(level protocol.LogLevel, message string) {
	h.conn.Notify(
		context.Background(),
		"window/logMessage",
		&protocol.LogMessageParams{
			Type:    level,
			Message: message,
		})
}
