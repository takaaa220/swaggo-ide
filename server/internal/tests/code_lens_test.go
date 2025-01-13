package tests

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/takaaa220/swaggo-ide/server/internal/handler/protocol"
)

func TestCodeLens(t *testing.T) {
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
},

// @Router / [post]
func main2() {
}

func main3() {
}
`,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res, err := sendTestRequest(t, testServer, "textDocument/codeLens", protocol.CodeLensParams{
		TextDocument: protocol.TextDocumentIdentifier{
			Uri: "file:///path/to/file.go",
		},
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedCommand := map[string]any{
		"title":   "swag fmt (this file)",
		"command": "swaggo-language-server-client.format",
		"arguments": []any{
			"/path/to/file.go",
		},
	}
	if diff := cmp.Diff([]any{
		map[string]any{
			"range": map[string]any{
				"start": map[string]any{
					"line":      float64(2),
					"character": float64(0),
				},
				"end": map[string]any{
					"line":      float64(2),
					"character": float64(0),
				},
			},
			"command": expectedCommand,
		},
		map[string]any{
			"range": map[string]any{
				"start": map[string]any{
					"line":      float64(8),
					"character": float64(0),
				},
				"end": map[string]any{
					"line":      float64(8),
					"character": float64(0),
				},
			},
			"command": expectedCommand,
		},
	}, res); diff != "" {
		t.Errorf("unexpected response (-want +got):\n%s", diff)
	}
}
