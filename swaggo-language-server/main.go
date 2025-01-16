package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal"
)

var debug = false

func main() {
	if debug {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	if err := internal.RunServer(debug); err != nil {
		if err != context.Canceled {
			log.Printf("failed to start server: %v", err)
		}
	}
}
