package protocol

import "github.com/takaaa220/go-swag-ide/server/v2/server-sdk/transport"

// see: https://www.jsonrpc.org/specification
// see: https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/

type InitializeFunc func(transport.Context, *InitializeParams) (*InitializeResult, error)

// InitializeParams defines the parameters sent in an initialize request.
type InitializeParams struct {
	ProcessID             int                `json:"processId,omitempty"`
	RootURI               string             `json:"rootUri,omitempty"`
	Capabilities          ClientCapabilities `json:"capabilities"`
	InitializationOptions interface{}        `json:"initializationOptions,omitempty"`
	Trace                 string             `json:"trace,omitempty"`
}

// InitializeResult defines the result returned for an initialize request.
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
}

// ClientCapabilities defines capabilities provided by the client.
type ClientCapabilities struct {
	TextDocument TextDocumentClientCapabilities `json:"textDocument,omitempty"`
	Workspace    WorkspaceClientCapabilities    `json:"workspace,omitempty"`
}

// ServerCapabilities defines the capabilities supported by the server.
type ServerCapabilities struct {
	TextDocumentSync   int                `json:"textDocumentSync"`
	HoverProvider      bool               `json:"hoverProvider,omitempty"`
	CompletionProvider *CompletionOptions `json:"completionProvider,omitempty"`
}
