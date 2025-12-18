package handler

import (
	"context"
	"time"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/handler/protocol"
	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/logger"
	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/swag"
)

type checkSyntaxRequest struct {
	uri  protocol.DocumentUri
	text string
}

func (h *LSPHandler) requestCheckSyntax(uri protocol.DocumentUri, text string) {
	h.checkSyntaxMu.Lock()
	defer h.checkSyntaxMu.Unlock()

	if h.checkSyntaxClosed {
		return
	}

	if h.checkSyntaxTimer != nil {
		h.checkSyntaxTimer.Stop()
		h.checkSyntaxTimer = nil
	}

	h.checkSyntaxTimer = time.AfterFunc(h.checkSyntaxDebounce, func() {
		h.checkSyntaxMu.Lock()
		h.checkSyntaxTimer = nil
		closed := h.checkSyntaxClosed
		h.checkSyntaxMu.Unlock()

		if closed {
			return
		}

		select {
		case h.checkSyntaxReq <- checkSyntaxRequest{
			uri:  uri,
			text: text,
		}:
		default:
			// channel is full or closed, skip
		}
	})
}

func (h *LSPHandler) checkSyntax(parentCtx context.Context) {
	runningCheckSyntax := make(map[protocol.DocumentUri]context.CancelFunc)

	for {
		select {
		case <-parentCtx.Done():
			logger.Infof("checkSyntax stopped")

			// Cancel all running syntax checks
			for uri, cancel := range runningCheckSyntax {
				cancel()
				delete(runningCheckSyntax, uri)
			}

			// Mark as closed and stop timer before closing channel
			h.checkSyntaxMu.Lock()
			h.checkSyntaxClosed = true
			if h.checkSyntaxTimer != nil {
				h.checkSyntaxTimer.Stop()
				h.checkSyntaxTimer = nil
			}
			h.checkSyntaxMu.Unlock()

			close(h.checkSyntaxReq)

			return
		case req, ok := <-h.checkSyntaxReq:
			if !ok {
				return
			}

			if cancel, ok := runningCheckSyntax[req.uri]; ok {
				cancel()
			}

			ctx, cancel := context.WithCancel(parentCtx)
			runningCheckSyntax[req.uri] = cancel

			go func(ctx context.Context, uri protocol.DocumentUri, text string) {
				defer cancel()

				syntaxErrors := swag.CheckSyntax(string(uri), text)

				// Skip notification if context is cancelled
				if ctx.Err() != nil {
					return
				}

				diagnostics := make([]protocol.Diagnostics, len(syntaxErrors))
				for i, syntaxError := range syntaxErrors {
					diagnostics[i] = protocol.Diagnostics{
						Range: protocol.Range{
							Start: protocol.Position{
								Line:      syntaxError.Range.Start.Line,
								Character: syntaxError.Range.Start.Character,
							},
							End: protocol.Position{
								Line:      syntaxError.Range.End.Line,
								Character: syntaxError.Range.End.Character,
							},
						},
						Severity: protocol.DiagnosticsSeverityError,
						Source:   "swag",
						Message:  syntaxError.Message,
					}
				}

				if err := h.Notify(ctx, "textDocument/publishDiagnostics",
					protocol.PublishDiagnosticsParams{
						Uri:         uri,
						Diagnostics: diagnostics,
					}); err != nil {
					if ctx.Err() == nil {
						logger.Error(err)
					}
				}
			}(ctx, req.uri, req.text)
		}
	}
}
