package tests

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/handler/protocol"
)

func TestInitialize(t *testing.T) {
	testServer := runServer()
	t.Cleanup(testServer.cancel)

	res, err := sendTestRequest(t, testServer, "initialize", protocol.InitializeParams{
		Capabilities: protocol.ClientCapabilities{},
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if diff := cmp.Diff(map[string]any{
		"capabilities": map[string]any{
			"completionProvider": map[string]any{
				"triggerCharacters": []any{" ", "@"},
			},
			"hoverProvider": true,
			"textDocumentSync": map[string]any{
				"openClose": true,
				"change":    1.0,
				"save": map[string]any{
					"includeText": true,
				},
			},
			"codeLensProvider": map[string]any{
				"resolveProvider": true,
			},
		},
	}, res); diff != "" {
		t.Errorf("unexpected response (-want +got):\n%s", diff)
	}
}
