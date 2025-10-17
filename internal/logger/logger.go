package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

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

// SetLevel changes the logging level
func SetLevel(level slog.Level) {
	handler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:      level,
		TimeFormat: time.Kitchen,
		NoColor:    !isTerminal(),
	})
	Log = slog.New(handler)
}

// isTerminal checks if stderr is a terminal
func isTerminal() bool {
	fileInfo, _ := os.Stderr.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
