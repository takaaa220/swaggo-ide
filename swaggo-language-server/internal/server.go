package internal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/handler"
	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/handler/transport"
	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/logger"
	"golang.org/x/exp/jsonrpc2"
)

func RunServer(config Config) error {
	server := newServer(config)
	if err := server.run(); err != nil && err != context.Canceled {
		return err
	}

	return nil
}

type server struct {
	config         Config
	jsonrpc2Server *jsonrpc2.Server
	cancel         context.CancelFunc
}

func newServer(config Config) *server {
	return &server{config: config}
}

func (s *server) run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.cancel = cancel

	var logWriter = os.Stderr
	if s.config.LogFile != "" {
		logFile, err := s.config.openLogFile()
		if err != nil {
			return err
		}
		defer logFile.Close()
		logWriter = logFile
	}
	logger.Setup(logWriter, s.config.LogLevel)

	if err := s.start(ctx); err != nil {
		return err
	}

	return s.wait(ctx)
}

func (s *server) start(ctx context.Context) error {
	handler := handler.NewLSPHandler(ctx, s.cancel)
	binder := transport.NewBinder(handler)
	listener := transport.NewStdListener()

	server, err := jsonrpc2.Serve(ctx, jsonrpc2.NewIdleListener(20*time.Second, listener), binder)
	if err != nil {
		return err
	}
	logger.Infof("Server started")

	s.jsonrpc2Server = server
	return nil
}

func (s *server) wait(ctx context.Context) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	go func() {
		select {
		case <-sigChan:
			s.cancel()
		case <-ctx.Done():
		}
	}()

	<-ctx.Done()

	// wait for graceful shutdown with timeout
	logger.Infof("Waiting for server to shut down...")

	done := make(chan struct{})
	go func() {
		s.jsonrpc2Server.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Infof("Server stopped")
	case <-time.After(5 * time.Second):
		logger.Infof("Server shutdown timed out, forcing exit")
	}

	return nil
}
