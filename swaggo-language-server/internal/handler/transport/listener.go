package transport

import (
	"context"
	"io"

	"golang.org/x/exp/jsonrpc2"
)

type stdListener struct {
	stdrwc *stdrwc
}

var _ jsonrpc2.Listener = (*stdListener)(nil)

func NewStdListener() jsonrpc2.Listener {
	return &stdListener{
		stdrwc: &stdrwc{},
	}
}

func (l *stdListener) Accept(ctx context.Context) (io.ReadWriteCloser, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return l.stdrwc, nil
	}
}

func (l *stdListener) Close() error {
	return l.stdrwc.Close()
}

func (l *stdListener) Dialer() jsonrpc2.Dialer {
	return nil
}
