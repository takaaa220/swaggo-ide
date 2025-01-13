package tests

import (
	"context"
	"errors"
	"log"
	"sync"
	"testing"

	"github.com/takaaa220/swaggo-ide/server/internal/handler"
	"github.com/takaaa220/swaggo-ide/server/internal/handler/transport"
	"golang.org/x/exp/jsonrpc2"
)

func sendTestRequest(t *testing.T, server *testServer, method string, params map[string]any) (any, error) {
	t.Helper()

	ctx := context.Background()

	conn, err := server.newConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	asyncCall := conn.Call(ctx, method, params)

	var result any
	if err := asyncCall.Await(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func runServer() (*testServer, func()) {
	ctx, cancel := context.WithCancel(context.Background())

	listener, err := jsonrpc2.NetPipe(ctx)
	if err != nil {
		log.Fatal(err)
	}

	binder := transport.NewBinder(handler.NewLSPHandler(handler.LSPHandlerOptions{LogLevel: handler.LogDebug}))

	server := &testServer{
		cancel:   cancel,
		listener: listener,
		binder:   binder,
	}

	s, err := jsonrpc2.Serve(
		ctx,
		listener,
		binder,
	)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := s.Wait(); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("Server stopped with error: %v", err)
			}
		}
	}()

	return server, cancel
}

type testServer struct {
	listener jsonrpc2.Listener
	binder   jsonrpc2.Binder
	cancel   func()
}

func (s *testServer) newConnection(ctx context.Context) (*jsonrpc2.Connection, error) {
	dialer := s.listener.Dialer()
	return jsonrpc2.Dial(ctx, dialer, s.binder)
}
