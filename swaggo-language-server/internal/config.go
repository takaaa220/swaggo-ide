package internal

import (
	"os"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/logger"
)

type Config struct {
	Debug    bool
	LogFile  string
	LogLevel logger.LogLevel
}

func (c *Config) openLogFile() (*os.File, error) {
	return os.OpenFile(c.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}
