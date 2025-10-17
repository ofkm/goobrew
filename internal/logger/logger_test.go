package logger

import (
	"log/slog"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	if Log == nil {
		t.Fatal("Log should be initialized")
	}
}

func TestSetLevel(t *testing.T) {
	tests := []struct {
		name  string
		level slog.Level
	}{
		{"Debug", slog.LevelDebug},
		{"Info", slog.LevelInfo},
		{"Warn", slog.LevelWarn},
		{"Error", slog.LevelError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLevel(tt.level)
			if Log == nil {
				t.Error("Log should not be nil after SetLevel")
			}
		})
	}
}

func TestIsTerminal(t *testing.T) {
	// Save original stderr
	oldStderr := os.Stderr
	defer func() { os.Stderr = oldStderr }()

	// Test with actual stderr
	result := isTerminal()
	// We can't assert the exact value as it depends on the test environment
	// Just ensure it doesn't panic
	_ = result
}

func TestLoggerWriting(t *testing.T) {
	// Test that we can write log messages without panicking
	SetLevel(slog.LevelDebug)

	Log.Debug("test debug message")
	Log.Info("test info message")
	Log.Warn("test warn message")
	Log.Error("test error message")

	Log.Debug("test with fields", "key", "value", "number", 42)
	Log.Info("test with multiple fields", "field1", "value1", "field2", "value2")
}

func TestSetLevelPersistence(t *testing.T) {
	originalLevel := slog.LevelInfo
	SetLevel(originalLevel)

	// Change level
	SetLevel(slog.LevelDebug)

	// Verify we can still log
	Log.Debug("should be visible at debug level")

	// Change back
	SetLevel(originalLevel)

	// Verify we can still log
	Log.Info("should be visible at info level")
}
