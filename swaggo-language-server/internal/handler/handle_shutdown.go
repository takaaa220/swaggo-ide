package handler

import (
	"context"

	"golang.org/x/exp/jsonrpc2"
)

func (h *LSPHandler) HandleShutdown(_ context.Context, req *jsonrpc2.Request) error {
	if h.checkSyntaxTimer != nil {
		h.checkSyntaxTimer.Stop()
	}

	h.logger.Debugf("Shutdown request received")
	h.shutdownChan <- struct{}{}
	h.CloseConnection()
	h.logger.Debugf("Shutdown completed")

	return nil
}
