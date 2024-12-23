package handler

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/takaaa220/swaggo-ide/server/internal/handler/filecache"
	"github.com/takaaa220/swaggo-ide/server/internal/handler/protocol"
	"github.com/takaaa220/swaggo-ide/server/internal/handler/swag"
	"golang.org/x/exp/jsonrpc2"
)

func NewLSPHandler(opts LSPHandlerOptions) *LSPHandler {
	logWriter := opts.logWriter
	if logWriter == nil {
		logWriter = os.Stderr
	}

	checkSyntaxDebounce := opts.checkSyntaxDebounce
	if checkSyntaxDebounce == 0 {
		checkSyntaxDebounce = 100 * time.Millisecond
	}

	h := &LSPHandler{
		checkSyntaxDebounce: checkSyntaxDebounce,
		checkSyntaxReq:      make(chan CheckSyntaxRequest),
		logger:              log.New(logWriter, "", log.LstdFlags),
	}

	go h.checkSyntax()
	return h
}

type LSPHandler struct {
	logger              *log.Logger
	conn                *jsonrpc2.Connection
	fileCache           *filecache.FileCache
	checkSyntaxReq      chan CheckSyntaxRequest
	checkSyntaxDebounce time.Duration
	checkSyntaxTimer    *time.Timer
}

type LSPHandlerOptions struct {
	checkSyntaxDebounce time.Duration
	logWriter           io.Writer
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
	default:
		h.logger.Printf("method %s not supported", req.Method)
		return nil, jsonrpc2.ErrNotHandled
	}
}

func (h *LSPHandler) SetConnection(conn *jsonrpc2.Connection) {
	// h.logger.Println("SetConnection called") // bug: this is called multiple times

	h.conn = conn
}

func (h *LSPHandler) CloseConnection() error {
	h.logger.Println("CloseConnection")

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
					h.logger.Println(err)
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
