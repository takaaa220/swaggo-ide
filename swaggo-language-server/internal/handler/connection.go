package handler

import "golang.org/x/exp/jsonrpc2"

func (h *LSPHandler) SetConnection(conn *jsonrpc2.Connection) {
	h.conn = conn
}

func (h *LSPHandler) CloseConnection() error {
	if h.conn == nil {
		return nil
	}
	return h.conn.Close()
}
