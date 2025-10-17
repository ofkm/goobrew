// Package cmd implements the command-line interface for goobrew.
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/ofkm/goobrew/internal/homebrew"
	"github.com/ofkm/goobrew/internal/logger"
	"github.com/ofkm/goobrew/internal/ui"
	"github.com/ofkm/goobrew/internal/version"
	"github.com/spf13/cobra"
)

// client is the Homebrew API client used across all commands.
var client *homebrew.Client

// verbose enables verbose output when set via the --verbose flag.
var verbose bool

// debug enables debug logging when set via the --debug flag.
var debug bool

// rootCmd represents the root command of goobrew.
var rootCmd = &cobra.Command{
	Use:   "goobrew",
	Short: "A fast and beautiful wrapper for Homebrew",
	Long: `goobrew is a modern Homebrew wrapper written in Go that uses 
Homebrew's JSON APIs for better performance and provides a 
beautiful, user-friendly interface.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set logging level
		if debug {
			logger.SetLevel(slog.LevelDebug)
		} else if verbose {
			logger.SetLevel(slog.LevelInfo)
		} else {
			logger.SetLevel(slog.LevelWarn)
		}

		// Initialize client
		var err error
		client, err = homebrew.NewClient()
		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		logger.Log.Debug("goobrew initialized", "version", version.Version)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		// Pass through to brew for unknown commands
		ctx := context.Background()
		if err := client.ExecuteCommand(ctx, args); err != nil {
			ui.PrintError(fmt.Sprintf("Command failed: %v", err))
			os.Exit(1)
		}
	},
}

// Execute runs the root command and all registered subcommands.
// This is the main entry point for command execution. It returns an error
// if command execution fails.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug output")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
