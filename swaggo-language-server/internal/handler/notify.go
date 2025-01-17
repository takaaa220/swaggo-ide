package handler

import "context"

func (h *LSPHandler) Notify(ctx context.Context, method string, params any) error {
	if h.conn == nil {
		return nil
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}

	return h.conn.Notify(ctx, method, params)
}
