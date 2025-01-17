package protocol

import "strings"

// file:///xxx/yyy/zzz
type DocumentUri string

func (d DocumentUri) Path() string {
	return strings.TrimPrefix(string(d), "file://")
}

type TextDocumentIdentifier struct {
	Uri DocumentUri `json:"uri"`
}

type Position struct {
	Line      uint32 `json:"line"`      // 0-based
	Character uint32 `json:"character"` // 0-based
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Location struct {
	Uri   DocumentUri `json:"uri"`
	Range Range       `json:"range"`
}

type ProgressToken any // string | int

type WorkDoneProgressParams struct {
	// An optional token that a server can use to report work done progress.
	WorkDoneToken ProgressToken `json:"workDoneToken,omitempty"`
}

type WorkDoneProgressOptions struct {
	WorkDoneProgress bool `json:"workDoneProgress,omitempty"`
}

type TextDocumentPositionParams struct {
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The position inside the text document.
	Position Position `json:"position"`
}

type MarkupKind string

const (
	// Plain text is supported as a content format
	MarkupKindPlainText MarkupKind = "plaintext"
	// Markdown is supported as a content format
	MarkupKindMarkdown MarkupKind = "markdown"
)

type MarkupContent struct {
	// The type of the Markup
	Kind MarkupKind `json:"kind"`
	// The content itself
	Value string `json:"value"`
}
