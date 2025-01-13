package protocol

type HoverClientCapabilities struct {
	// Whether hover supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// Client supports the follow content formats if the content
	// property refers to a `literal of type MarkupContent`.
	// The order describes the preferred format of the client.
	ContentFormat []MarkupKind `json:"contentFormat,omitempty"`
}

type MarkupKind string

const (
	// Plain text is supported as a content format
	MarkupKindPlainText MarkupKind = "plaintext"
	// Markdown is supported as a content format
	MarkupKindMarkdown MarkupKind = "markdown"
)

type HoverOptions struct {
	WorkDoneProgressOptions
}

type HoverRegistrationOptions struct {
	TextDocumentRegistrationOptions
	HoverOptions
}

type HoverParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
}

type Hover struct {
	// The hover's content
	Contents MarkupContent `json:"contents"`
	// An optional range is a range inside a text document
	// that is used to visualize a hover, e.g. by changing the background color.
	Range *Range `json:"range,omitempty"`
}

type MarkupContent struct {
	// The type of the Markup
	Kind MarkupKind `json:"kind"`
	// The content itself
	Value string `json:"value"`
}
