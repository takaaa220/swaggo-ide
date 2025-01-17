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
	if h.checkSyntaxTimer != nil {
		h.checkSyntaxTimer.Reset(h.checkSyntaxDebounce)
	}

	h.checkSyntaxTimer = time.AfterFunc(h.checkSyntaxDebounce, func() {
		h.checkSyntaxTimer = nil
		h.checkSyntaxReq <- checkSyntaxRequest{
			uri:  uri,
			text: text,
		}
	})
}

func (h *LSPHandler) checkSyntax(ctx context.Context) {
	runningCheckSyntax := make(map[protocol.DocumentUri]context.CancelFunc)

	for {
		select {
		case <-ctx.Done():
			logger.Infof("checkSyntax stopped")
			// FIXME: refactor (shouldn't stop timer here)
			for uri, cancel := range runningCheckSyntax {
				cancel()
				delete(runningCheckSyntax, uri)
			}
			if h.checkSyntaxTimer != nil {
				h.checkSyntaxTimer.Stop()
			}
			close(h.checkSyntaxReq)

			return
		case req, ok := <-h.checkSyntaxReq:
			if !ok {
				break
			}

			if cancel, ok := runningCheckSyntax[req.uri]; ok {
				cancel()
			}

			ctx, cancel := context.WithCancel(context.Background())
			runningCheckSyntax[req.uri] = cancel

			go func(uri protocol.DocumentUri, text string) {
				defer cancel()

				syntaxErrors := swag.CheckSyntax(string(uri), text)

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
					logger.Error(err)
				}
			}(req.uri, req.text)
		}
	}
}
