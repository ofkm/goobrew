package main

import (
	"os"

	"github.com/ofkm/goobrew/cmd"
	"github.com/ofkm/goobrew/internal/logger"
	"github.com/ofkm/goobrew/internal/ui"
)

func main() {
	if err := cmd.Execute(); err != nil {
		ui.PrintError(err.Error())
		logger.Log.Error("command execution failed", "error", err)
		os.Exit(1)
	}
}
