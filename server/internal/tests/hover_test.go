package tests

import (
	"fmt"
	"testing"
)

func TestHover(t *testing.T) {
	t.Parallel()

	t.Run("TestHover", func(t *testing.T) {
		testServer, cancel := runServer()
		defer cancel()

		res, err := sendTestRequest(t, testServer, "initialize", map[string]any{
			"hello": "world",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		fmt.Println(res)

		res, err = sendTestRequest(t, testServer, "textDocument/hover", map[string]any{
			"hello": "world",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		fmt.Println(res)
	})

	t.Run("TestHover2", func(t *testing.T) {
		testServer, cancel := runServer()
		defer cancel()

		res, err := sendTestRequest(t, testServer, "initialize", map[string]any{
			"hello": "world",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		fmt.Println(res)

		res, err = sendTestRequest(t, testServer, "textDocument/hover", map[string]any{
			"hello": "world",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		fmt.Println(res)
	})

}
