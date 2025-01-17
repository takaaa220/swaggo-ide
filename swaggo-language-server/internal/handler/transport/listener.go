package transport

import (
	"context"
	"io"
	"sync"

	"golang.org/x/exp/jsonrpc2"
)

type stdListener struct {
	stdrwc io.ReadWriteCloser
	mu     sync.Mutex
}

func NewStdListener() *stdListener {
	return &stdListener{}
}

func (l *stdListener) Accept(ctx context.Context) (io.ReadWriteCloser, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	l.stdrwc = &stdrwc{}
	return l.stdrwc, nil
}

func (l *stdListener) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.stdrwc == nil {
		return nil
	}
	if err := l.stdrwc.Close(); err != nil {
		return err
	}

	l.stdrwc = nil
	return nil
}

func (l *stdListener) Dialer() jsonrpc2.Dialer {
	return nil
}
