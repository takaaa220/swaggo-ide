package handler

import (
	"context"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/logger"
	"golang.org/x/exp/jsonrpc2"
)

func (h *LSPHandler) HandleShutdown(_ context.Context, req *jsonrpc2.Request) error {
	if h.checkSyntaxTimer != nil {
		h.checkSyntaxTimer.Stop()
	}

	logger.Debugf("Shutdown request received")
	h.cancel()
	h.CloseConnection()
	logger.Debugf("Shutdown completed")

	return nil
}
