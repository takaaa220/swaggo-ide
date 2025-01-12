package protocol

type CodeLensClientCapabilities struct {
	// Whether code lens supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

type CodeLensOptions struct {
	WorkDoneProgress bool `json:"workDoneProgress,omitempty"`
	// Code lens has a resolve provider as well.
	ResolveProvider bool `json:"resolveProvider,omitempty"`
}

type CodeLensParams struct {
	// The document to request code lens for.
	TextDocument       TextDocumentIdentifier `json:"textDocument"`
	WorkDoneToken      ProgressToken          `json:"workDoneToken,omitempty"`
	PartialResultToken ProgressToken          `json:"partialResultToken,omitempty"`
}

type CodeLens struct {
	// The range in which this code lens is valid. Should only span a single line.
	Range Range `json:"range"`
	// The command this code lens represents.
	Command Command `json:"command,omitempty"`
	// A data entry field that is preserved on a code lens item between
	// a code lens and a code lens resolve request.
	Data any `json:"data,omitempty"`
}

type Command struct {
	// Title of the command, like `save`.
	Title string `json:"title"`
	// The identifier of the actual command handler.
	Command string `json:"command"`
	// Arguments that the command handler should be
	// invoked with.
	Arguments []any `json:"arguments,omitempty"`
}
