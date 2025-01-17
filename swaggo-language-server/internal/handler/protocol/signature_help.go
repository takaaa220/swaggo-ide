package protocol

type SignatureHelpClientCapabilities struct {
	// Whether signature help supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// The client supports the following `SignatureInformation` specific properties.
	SignatureInformation *SignatureInformationInner `json:"signatureInformation,omitempty"`

	// The client supports to send additional context information for a `textDocument/signatureHelp` request.
	// A client that opts into contextSupport will also support the `retriggerCharacters` on `SignatureHelpOptions`.
	// @since 3.15.0
	ContextSupport bool `json:"contextSupport,omitempty"`
}

type SignatureInformationInner struct {
	// Client supports the follow content formats for the documentation property.
	// The order describes the preferred format of the client.
	DocumentationFormat []MarkupKind `json:"documentationFormat,omitempty"`

	// Client capabilities specific to parameter information.
	ParameterInformation *ParameterInformationInner `json:"parameterInformation,omitempty"`

	// The client supports the `activeParameter` property on `SignatureInformation` literal.
	// @since 3.16.0
	ActiveParameterSupport bool `json:"activeParameterSupport,omitempty"`
}

type ParameterInformationInner struct {
	// The client supports processing label offsets instead of a simple label string.
	// @since 3.14.0
	LabelOffsetSupport bool `json:"labelOffsetSupport,omitempty"`
}

// SignatureHelpOptions represents the options for signature help in the Language Server Protocol.
type SignatureHelpOptions struct {
	WorkDoneProgressOptions

	// The characters that trigger signature help automatically.
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`

	// List of characters that re-trigger signature help.
	//
	// These trigger characters are only active when signature help is already
	// showing. All trigger characters are also counted as re-trigger
	// characters.
	//
	// @since 3.15.0
	RetriggerCharacters []string `json:"retriggerCharacters,omitempty"`
}

// SignatureHelpParams represents the parameters for a signature help request.
type SignatureHelpParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams

	// The signature help context. This is only available if the client
	// specifies to send this using the client capability
	// `textDocument.signatureHelp.contextSupport === true`
	//
	// @since 3.15.0
	Context *SignatureHelpContext `json:"context,omitempty"`
}

// SignatureHelpTriggerKind represents how a signature help was triggered.
//
// @since 3.15.0
type SignatureHelpTriggerKind int

const (
	// Signature help was invoked manually by the user or by a command.
	SignatureHelpTriggerKindInvoked SignatureHelpTriggerKind = 1
	// Signature help was triggered by a trigger character.
	SignatureHelpTriggerKindTriggerCharacter SignatureHelpTriggerKind = 2
	// Signature help was triggered by the cursor moving or by the document
	// content changing.
	SignatureHelpTriggerKindContentChange SignatureHelpTriggerKind = 3
)

// SignatureHelpContext provides additional information about the context in which a signature help request was triggered.
//
// @since 3.15.0
type SignatureHelpContext struct {
	// Action that caused signature help to be triggered.
	TriggerKind SignatureHelpTriggerKind `json:"triggerKind"`

	// Character that caused signature help to be triggered.
	//
	// This is undefined when triggerKind !== SignatureHelpTriggerKindTriggerCharacter
	TriggerCharacter string `json:"triggerCharacter,omitempty"`

	// `true` if signature help was already showing when it was triggered.
	//
	// Retriggers occur when the signature help is already active and can be
	// caused by actions such as typing a trigger character, a cursor move, or
	// document content changes.
	IsRetrigger bool `json:"isRetrigger"`

	// The currently active `SignatureHelp`.
	//
	// The `activeSignatureHelp` has its `SignatureHelp.activeSignature` field
	// updated based on the user navigating through available signatures.
	ActiveSignatureHelp *SignatureHelp `json:"activeSignatureHelp,omitempty"`
}

// SignatureHelp represents the signature of something callable.
// There can be multiple signatures but only one active and only one active parameter.
type SignatureHelp struct {
	// Signatures is one or more signatures. If no signatures are available, the signature help
	// request should return nil.
	Signatures []SignatureInformation `json:"signatures"`

	// ActiveSignature is the active signature. If omitted or the value lies outside the
	// range of signatures, the value defaults to zero or is ignored if the
	// SignatureHelp has no signatures.
	//
	// Whenever possible, implementors should make an active decision about
	// the active signature and shouldn't rely on a default value.
	//
	// In future versions of the protocol, this property might become
	// mandatory to better express this.
	ActiveSignature *uint `json:"activeSignature,omitempty"`

	// ActiveParameter is the active parameter of the active signature. If omitted or the value
	// lies outside the range of signatures[activeSignature].parameters,
	// defaults to 0 if the active signature has parameters. If
	// the active signature has no parameters, it is ignored.
	//
	// In future versions of the protocol, this property might become
	// mandatory to better express the active parameter if the
	// active signature does have any.
	ActiveParameter *uint `json:"activeParameter,omitempty"`
}

// SignatureInformation represents the signature of something callable.
// A signature can have a label, like a function name, a doc-comment, and
// a set of parameters.
type SignatureInformation struct {
	// Label is the label of this signature. Will be shown in the UI.
	Label string `json:"label"`

	// Documentation is the human-readable doc-comment of this signature. Will be shown
	// in the UI but can be omitted.
	Documentation *MarkupContent `json:"documentation,omitempty"`

	// Parameters is the list of parameters of this signature.
	Parameters []ParameterInformation `json:"parameters,omitempty"`

	// ActiveParameter is the index of the active parameter.
	//
	// If provided, this is used in place of SignatureHelp.ActiveParameter.
	//
	// Available since version 3.16.0.
	ActiveParameter *uint `json:"activeParameter,omitempty"`
}

// ParameterInformation represents a parameter of a callable signature.
// A parameter can have a label and a doc-comment.
type ParameterInformation struct {
	// Label is the label of this parameter information.
	//
	// Either a string or an inclusive start and exclusive end offsets within
	// its containing signature label (see SignatureInformation.Label). The
	// offsets are based on a UTF-16 string representation as Position and
	// Range do.
	//
	// Note: a label of type string should be a substring of its containing
	// signature label. Its intended use case is to highlight the parameter
	// label part in the SignatureInformation.Label.
	Label string `json:"label"` // Can be string or [2]uint

	// Documentation is the human-readable doc-comment of this parameter. Will be shown
	// in the UI but can be omitted.
	Documentation *MarkupContent `json:"documentation,omitempty"`
}
