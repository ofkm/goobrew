package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/ofkm/goobrew/internal/homebrew"
	"github.com/ofkm/goobrew/internal/logger"
	"github.com/ofkm/goobrew/internal/ui"
	"github.com/spf13/cobra"
)

var (
	client  *homebrew.Client
	verbose bool
	debug   bool
)

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

		logger.Log.Debug("goobrew initialized", "version", "dev")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
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

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug output")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
