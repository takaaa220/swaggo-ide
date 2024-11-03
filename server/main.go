package main

import (
	_ "net/http/pprof"

	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/takaaa220/go-swag-ide/server/v2/server"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	ctx := context.Background()

	if err := server.StartServer(ctx); err != nil {
		log.Fatal(fmt.Errorf("failed to start server: %w", err))
	}
}
