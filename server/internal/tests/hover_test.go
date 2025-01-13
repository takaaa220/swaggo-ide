package tests

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/takaaa220/swaggo-ide/server/internal/handler/protocol"
)

func TestHover(t *testing.T) {
	testServer := runServer()
	t.Cleanup(testServer.cancel)

	if _, err := sendTestRequest(t, testServer, "initialize", protocol.InitializeParams{
		Capabilities: protocol.ClientCapabilities{},
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := sendTestRequest(t, testServer, "textDocument/didSave", protocol.DidSaveTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{
			Uri: "file:///path/to/file.go",
		},
		Text: `package main

// main is the entry point for the program.
// @Router / [get]
// @Summary Get the main page
func main() {
}`,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res, err := sendTestRequest(t, testServer, "textDocument/hover", protocol.HoverParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				Uri: "file:///path/to/file.go",
			},
			Position: protocol.Position{
				Line:      3,
				Character: 4,
			},
		},
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if diff := cmp.Diff(map[string]any{
		"contents": map[string]any{
			"kind":  "markdown",
			"value": "```\n@Router PATH [HTTP_METHOD]\n```\n\n---\nA router definition. This is required for the operation.",
		},
	}, res); diff != "" {
		t.Errorf("unexpected response (-want +got):\n%s", diff)
	}
}
