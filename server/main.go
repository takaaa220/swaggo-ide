package main

import (
	"fmt"

	"github.com/takaaa220/go-swag-ide/server/internal"
)

func main() {
	if err := internal.StartServer(); err != nil {
		panic(fmt.Errorf("failed to start server: %w", err))
	}
}
