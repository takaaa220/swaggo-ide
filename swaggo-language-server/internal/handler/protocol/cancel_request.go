package protocol

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/jsonrpc2"
)

type CancelParams struct {
	ID jsonrpc2.ID `json:"id"` // string | int64
}

func (c *CancelParams) UnmarshalJSON(b []byte) error {
	p := struct {
		ID any `json:"id"`
	}{}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	switch v := p.ID.(type) {
	case string:
		c.ID = jsonrpc2.StringID(v)
	case float64:
		c.ID = jsonrpc2.Int64ID(int64(v))
	default:
		return fmt.Errorf("invalid message id type <%T>%v", v, v)
	}

	return nil
}
