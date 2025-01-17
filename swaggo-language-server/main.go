package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal"
	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/logger"
)

func main() {
	cfg := config()

	if cfg.Debug {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	if err := internal.RunServer(cfg); err != nil {
		logger.Errorf("failed to run server: %v", err)
	}
}

func config() internal.Config {
	// TODO: parse command line arguments
	debug := false

	if debug {
		return internal.Config{
			LogFile:  "swaggo-language-server.log",
			LogLevel: logger.LogDebug,
		}
	}

	return internal.Config{
		LogFile:  "",
		LogLevel: logger.LogInfo,
	}
}
