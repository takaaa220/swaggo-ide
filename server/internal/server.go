package internal

import (
	"log"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const lsName = "Language Server for Go Swag"

var (
	version string = "0.0.1"
	handler protocol.Handler
)

func StartServer() error {

	handler := protocol.Handler{
		Initialize:                     initialize,
		Initialized:                    initialized,
		Shutdown:                       shutdown,
		SetTrace:                       setTrace,
		TextDocumentDidSave:            didSave,
		WorkspaceDidChangeWatchedFiles: didChangeWatchedFiles,
	}

	s := server.NewServer(&handler, lsName, false)
	return s.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	log.Println("Initializing Language Server")

	capabilities := handler.CreateServerCapabilities()

	capabilities.TextDocumentSync = protocol.TextDocumentSyncOptions{
		Save: true,
	}
	capabilities.CodeLensProvider = &protocol.CodeLensOptions{}

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func didSave(context *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
	log.Println("Saved", params.TextDocument.URI)

	go context.Notify(protocol.ServerWorkspaceCodeLensRefresh, nil)

	return nil
}

func didChangeWatchedFiles(context *glsp.Context, params *protocol.DidChangeWatchedFilesParams) error {
	log.Println("Changed watched files")

	return nil
}
