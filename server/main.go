package main

import (
	"context"
	"fmt"
	"log"

	"github.com/takaaa220/go-swag-ide/server/v2/server"
)

func main() {
	ctx := context.Background()

	if err := server.StartServer(ctx); err != nil {
		log.Fatal(fmt.Errorf("failed to start server: %w", err))
	}
}
