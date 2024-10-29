package protocol

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
