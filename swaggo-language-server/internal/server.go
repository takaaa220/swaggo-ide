package internal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/handler"
	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/handler/transport"
	"golang.org/x/exp/jsonrpc2"
)

func RunServer(debug bool) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: refactor for options
	opts := handler.LSPHandlerOptions{
		CheckSyntaxDebounce: 100 * time.Millisecond,
		LogLevel:            handler.LogWarn,
		LogWriter:           os.Stderr,
	}
	if debug {
		opts.LogLevel = handler.LogDebug

		logWriter, err := os.OpenFile("swaggo-language-server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer logWriter.Close()
		opts.LogWriter = logWriter
	}

	shutdownChan := make(chan struct{})
	defer close(shutdownChan)

	handler := handler.NewLSPHandler(ctx, shutdownChan, opts)
	binder := transport.NewBinder(handler)
	listener := transport.NewStdListener()

	server, err := jsonrpc2.Serve(ctx, jsonrpc2.NewIdleListener(20*time.Second, listener), binder)
	if err != nil {
		return err
	}

	fmt.Fprintln(opts.LogWriter, "Server started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	select {
	case <-sigChan:
	case <-shutdownChan:
		cancel()
	case <-ctx.Done():
	}

	// wait for graceful shutdown
	fmt.Fprintln(opts.LogWriter, "Waiting for server to shut down...")
	server.Wait()
	fmt.Fprintln(opts.LogWriter, "Server stopped")
	return nil
}
