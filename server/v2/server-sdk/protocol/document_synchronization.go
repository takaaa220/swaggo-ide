package protocol

import "github.com/takaaa220/go-swag-ide/server/v2/server-sdk/transport"

// TextDocumentClientCapabilities defines capabilities related to text documents.
type TextDocumentClientCapabilities struct {
	Synchronization TextDocumentSyncClientCapabilities `json:"synchronization,omitempty"`
}

// TextDocumentSyncClientCapabilities defines capabilities for text document synchronization.
type TextDocumentSyncClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	WillSave            bool `json:"willSave,omitempty"`
	WillSaveWaitUntil   bool `json:"willSaveWaitUntil,omitempty"`
	DidSave             bool `json:"didSave,omitempty"`
}

type TextDocumentServerCapabilities struct {
	TextDocumentSync  any  `json:"textDocumentSync,omitempty"` // TextDocumentSyncOptions | TextDocumentSyncKind
	Save              any  `json:"save,omitempty"`             // SaveOptions | bool
	WillSave          bool `json:"willSave,omitempty"`
	WillSaveWaitUntil bool `json:"willSaveWaitUntil,omitempty"`
}

type TextDocumentSyncKind int

const (
	TextDocumentSyncKindNone        TextDocumentSyncKind = 0
	TextDocumentSyncKindFull        TextDocumentSyncKind = 1
	TextDocumentSyncKindIncremental TextDocumentSyncKind = 2
)

type TextDocumentSyncOptions struct {
	OpenClose bool                 `json:"openClose,omitempty"`
	Change    TextDocumentSyncKind `json:"change,omitempty"`
}

type DidOpenTextDocumentFunc func(context transport.Context, params *DidOpenTextDocumentParams) error

type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

type TextDocumentItem struct {
	Uri        DocumentUri `json:"uri"`
	LanguageId string      `json:"languageId"`
	Version    int         `json:"version"`
	Text       string      `json:"text"`
}

type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type DidChangeTextDocumentFunc func(context transport.Context, params *DidChangeTextDocumentParams) error

type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}

type TextDocumentContentChangeEvent struct {
	Range       Range  `json:"range,omitempty"`
	RangeLength int    `json:"rangeLength,omitempty"`
	Text        string `json:"text"`
}

type TextDocumentRegistrationOptions struct {
	DocumentSelector DocumentSelector `json:"documentSelector,omitempty"`
}

type DocumentFilter struct {
	Language string `json:"language,omitempty"`
	Scheme   string `json:"scheme,omitempty"`
	Pattern  string `json:"pattern,omitempty"`
}

type DocumentSelector []DocumentFilter

type TextDocumentChangeRegistrationOptions struct {
	TextDocumentRegistrationOptions

	SyncKind TextDocumentSyncKind `json:"syncKind,omitempty"`
}

type WillSaveTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Reason       TextDocumentSaveReason `json:"reason"`
}

type TextDocumentSaveReason int

const (
	TextDocumentSaveReasonManual     TextDocumentSaveReason = 1
	TextDocumentSaveReasonAfterDelay TextDocumentSaveReason = 2
	TextDocumentSaveReasonFocusOut   TextDocumentSaveReason = 3
)

type SaveOptions struct {
	IncludeText bool `json:"includeText,omitempty"`
}

type TextDocumentSaveRegistrationOptions struct {
	TextDocumentRegistrationOptions

	IncludeText bool `json:"includeText,omitempty"`
}

type TextDocumentDidSaveFunc func(context transport.Context, params *DidSaveTextDocumentParams) error

type DidSaveTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Text         string                 `json:"text,omitempty"`
}

type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type TextDocumentDidCloseFunc func(context transport.Context, params *DidCloseTextDocumentParams) error
