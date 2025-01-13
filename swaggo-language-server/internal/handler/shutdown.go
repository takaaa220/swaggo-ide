package handler

import (
	"context"

	"golang.org/x/exp/jsonrpc2"
)

func (h *LSPHandler) HandleShutdown(_ context.Context, _ *jsonrpc2.Request) (any, error) {
	if h.checkSyntaxTimer != nil {
		h.checkSyntaxTimer.Stop()
	}

	close(h.checkSyntaxReq)
	return nil, h.CloseConnection()
}
