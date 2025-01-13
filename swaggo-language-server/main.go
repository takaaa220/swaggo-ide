package main

import (
	_ "net/http/pprof"

	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal"
)

var debug = false

func main() {
	if debug {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	ctx := context.Background()

	if err := internal.StartServer(ctx, debug); err != nil {
		log.Fatal(fmt.Errorf("failed to start server: %w", err))
	}
}
