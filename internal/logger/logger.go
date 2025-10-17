// Package logger provides structured logging functionality for goobrew.
// It uses the slog package with tint for colorized output to stderr.
package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

// Log is the global logger instance used throughout the application.
// It is initialized with tint handler for colorized output and defaults
// to Info level logging.
var Log *slog.Logger

func init() {
	// Create a tint handler with custom options
	handler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.Kitchen,
		NoColor:    !isTerminal(),
	})

	Log = slog.New(handler)
}

// SetLevel changes the logging level for the global logger.
// It recreates the logger with a new tint handler configured for the specified level.
// Valid levels are slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, and slog.LevelError.
func SetLevel(level slog.Level) {
	handler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:      level,
		TimeFormat: time.Kitchen,
		NoColor:    !isTerminal(),
	})
	Log = slog.New(handler)
}

// isTerminal checks if stderr is connected to a terminal.
// It returns true if stderr is a character device (terminal), false otherwise.
// This is used to determine whether to enable colored output.
func isTerminal() bool {
	fileInfo, _ := os.Stderr.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
