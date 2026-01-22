package test

import (
	"testing"

	"log/slog"

	"github.com/leijux/templog"
)

func TestTemplog(t *testing.T) {
	defer templog.Close()

	slog.Info("This is an info message")
	slog.Debug("This is a debug message")
	slog.Warn("This is a warning message")
	slog.Error("This is an error message")
}
