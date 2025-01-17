package handler

import (
	"context"
	"io"
	"time"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/handler/filecache"
	"golang.org/x/exp/jsonrpc2"
)

func NewLSPHandler(ctx context.Context, shutdownChan chan struct{}, opts LSPHandlerOptions) *LSPHandler {
	h := &LSPHandler{
		checkSyntaxDebounce: opts.CheckSyntaxDebounce,
		checkSyntaxReq:      make(chan checkSyntaxRequest),
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
	checkSyntaxReq      chan checkSyntaxRequest
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

// this is hack (error occurs when Handle returns (nil, nil))
type Null struct{}

func (n *Null) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}
