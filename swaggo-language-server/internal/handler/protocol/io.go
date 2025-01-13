package protocol

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/jsonrpc2"
)

// Error codes as per the JSON-RPC 2.0 specification and LSP custom errors.
type LSPCodeError int

const (
	CodeParseError     LSPCodeError = -32700
	CodeInvalidRequest              = -32600
	CodeMethodNotFound              = -32601
	CodeInvalidParams               = -32602
	CodeInternalError               = -32603
	CodeServerError                 = -32000
)

func NewSuccessResponse(id uint64, result any) (*jsonrpc2.Response, error) {
	b, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}
	rawResult := json.RawMessage(b)

	return &jsonrpc2.Response{
		ID:     jsonrpc2.Int64ID(int64(id)),
		Result: rawResult,
	}, nil
}

func NewFailureResponse(id uint64, code LSPCodeError, message string, data any) (*jsonrpc2.Response, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal error data: %w", err)
	}
	rawData := json.RawMessage(b)

	responseError := NewResponseError(code, message, rawData)
	return &jsonrpc2.Response{
		ID:    jsonrpc2.Int64ID(int64(id)),
		Error: responseError,
	}, nil
}

type Notification struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}

// NewNotification creates a new Notification with a specific method and params.
func NewNotification(method string, params any) (*jsonrpc2.Request, error) {
	return jsonrpc2.NewNotification(method, params)
}
