package repository_test

import (
	"os"
	"testing"
	"yotudo/src/lib/logger"
)

func TestMain(m *testing.M) {
	logger.InitializeLogger("debug", []string{logger.Console_Type})

	if c := m.Run(); c != 0 {
		os.Exit(c)
	}
}
