package tests

import (
	"os"
	"testing"

	"github.com/takaaa220/swaggo-ide/swaggo-language-server/internal/logger"
)

func TestMain(m *testing.M) {
	logger.Setup(os.Stderr, logger.LogDebug)

	os.Exit(m.Run())
}
