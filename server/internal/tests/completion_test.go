package tests

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/takaaa220/swaggo-ide/server/internal/handler/protocol"
)

func TestCompletion(t *testing.T) {
	t.Run("complete api operations", func(t *testing.T) {
		testServer := runServer()
		t.Cleanup(testServer.cancel)

		if _, err := sendTestRequest(t, testServer, "initialize", protocol.InitializeParams{
			Capabilities: protocol.ClientCapabilities{},
		}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, err := sendTestRequest(t, testServer, "textDocument/didOpen", protocol.DidOpenTextDocumentParams{
			TextDocument: protocol.TextDocumentItem{
				Uri:        "file:///path/to/file.go",
				LanguageId: "go",
				Version:    1,
				Text: `package main

// main is the entry point for the program.
// @Router / [get]
// @Summary Get the main page
func main() {
}`,
			},
		}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, err := sendTestRequest(t, testServer, "textDocument/didChange", protocol.DidChangeTextDocumentParams{
			TextDocument: protocol.VersionedTextDocumentIdentifier{
				TextDocumentIdentifier: protocol.TextDocumentIdentifier{
					Uri: "file:///path/to/file.go",
				},
				Version: 2,
			},
			ContentChanges: []protocol.TextDocumentContentChangeEvent{
				{
					Text: `package main

// main is the entry point for the program.
// @Router / [get]
// @Summary Get the main page
// @P
func main() {
}`,
				},
			},
		}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// wait for setting cache
		time.Sleep(100 * time.Millisecond)

		res, err := sendTestRequest(t, testServer, "textDocument/completion", protocol.CompletionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				Uri: "file:///path/to/file.go",
			},
			Position: protocol.Position{
				Line:      5,
				Character: 4,
			},
			Context: protocol.CompletionContext{
				TriggerKind:      protocol.CompletionTriggerKindCharacter,
				TriggerCharacter: "@",
			},
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expectedCompletionItem := func(label, newText string) map[string]any {
			return map[string]any{
				"kind":  float64(14),
				"label": label,
				"textEdit": map[string]any{
					"newText": newText,
					"range": map[string]any{
						"start": map[string]any{
							"line":      float64(5),
							"character": float64(3),
						},
						"end": map[string]any{
							"line":      float64(5),
							"character": float64(4),
						},
					},
				},
			}
		}
		if diff := cmp.Diff(map[string]any{
			"isIncomplete": false,
			"items": []any{
				expectedCompletionItem("@Accept", "@Accept MIME_TYPE"),
				expectedCompletionItem("@Description", "@Description DESCRIPTION"),
				expectedCompletionItem("@Failure", "@Failure STATUS_CODE {DATA_TYPE} GO_TYPE"),
				expectedCompletionItem("@Header", "@Header STATUS_CODE {DATA_TYPE} HEADER_NAME COMMENT"),
				expectedCompletionItem("@ID", "@ID ID"),
				expectedCompletionItem("@Param", "@Param PARAM_NAME PARAM_TYPE GO_TYPE REQUIRED \"DESCRIPTION\""),
				expectedCompletionItem("@Produce", "@Produce MIME_TYPE"),
				expectedCompletionItem("@Router", "@Router PATH [HTTP_METHOD]"),
				expectedCompletionItem("@Success", "@Success STATUS_CODE {DATA_TYPE} GO_TYPE"),
				expectedCompletionItem("@Summary", "@Summary SUMMARY"),
				expectedCompletionItem("@Tags", "@Tags TAG1,TAG2"),
			},
		}, res); diff != "" {
			t.Errorf("unexpected response (-want +got):\n%s", diff)
		}
	})

	t.Run("complete tag args", func(t *testing.T) {
		testServer := runServer()
		t.Cleanup(testServer.cancel)

		if _, err := sendTestRequest(t, testServer, "initialize", protocol.InitializeParams{
			Capabilities: protocol.ClientCapabilities{},
		}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, err := sendTestRequest(t, testServer, "textDocument/didOpen", protocol.DidOpenTextDocumentParams{
			TextDocument: protocol.TextDocumentItem{
				Uri:        "file:///path/to/file.go",
				LanguageId: "go",
				Version:    1,
				Text: `package main

// main is the entry point for the program.
// @Router / [get]
// @Summary Get the main page
func main() {
}`,
			},
		}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, err := sendTestRequest(t, testServer, "textDocument/didChange", protocol.DidChangeTextDocumentParams{
			TextDocument: protocol.VersionedTextDocumentIdentifier{
				TextDocumentIdentifier: protocol.TextDocumentIdentifier{
					Uri: "file:///path/to/file.go",
				},
				Version: 2,
			},
			ContentChanges: []protocol.TextDocumentContentChangeEvent{
				{
					Text: `package main

// main is the entry point for the program.
// @Router / [get]
// @Summary Get the main page
// @Param id p
func main() {
}`,
				},
			},
		}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// wait for setting cache
		time.Sleep(100 * time.Millisecond)

		res, err := sendTestRequest(t, testServer, "textDocument/completion", protocol.CompletionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				Uri: "file:///path/to/file.go",
			},
			Position: protocol.Position{
				Line:      5,
				Character: 14,
			},
			Context: protocol.CompletionContext{
				TriggerKind:      protocol.CompletionTriggerKindCharacter,
				TriggerCharacter: " ",
			},
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expectedCompletionItem := func(label, newText string) map[string]any {
			return map[string]any{
				"kind":  float64(14),
				"label": label,
				"textEdit": map[string]any{
					"newText": newText,
					"range": map[string]any{
						"start": map[string]any{
							"line":      float64(5),
							"character": float64(14),
						},
						"end": map[string]any{
							"line":      float64(5),
							"character": float64(14),
						},
					},
				},
			}
		}
		if diff := cmp.Diff(map[string]any{
			"isIncomplete": false,
			"items": []any{
				expectedCompletionItem("string", "string"),
				expectedCompletionItem("number", "number"),
				expectedCompletionItem("integer", "integer"),
				expectedCompletionItem("boolean", "boolean"),
				expectedCompletionItem("file", "file"),
				expectedCompletionItem("object", "object"),
			},
		}, res); diff != "" {
			t.Errorf("unexpected response (-want +got):\n%s", diff)
		}
	})
}
