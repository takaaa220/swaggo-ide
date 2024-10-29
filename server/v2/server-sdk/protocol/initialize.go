package protocol

// see: https://www.jsonrpc.org/specification
// see: https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/

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

// CompletionOptions represents options for completion requests.
type CompletionOptions struct {
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
	ResolveProvider   bool     `json:"resolveProvider,omitempty"`
}

// TextDocumentClientCapabilities defines capabilities related to text documents.
type TextDocumentClientCapabilities struct {
	Synchronization TextDocumentSyncClientCapabilities `json:"synchronization,omitempty"`
}

// WorkspaceClientCapabilities defines capabilities related to the workspace.
type WorkspaceClientCapabilities struct {
	ApplyEdit bool `json:"applyEdit,omitempty"`
}

// TextDocumentSyncClientCapabilities defines capabilities for text document synchronization.
type TextDocumentSyncClientCapabilities struct {
	DidSave bool `json:"didSave,omitempty"`
}
