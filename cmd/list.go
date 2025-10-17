package cmd

import (
	"context"
	"os"

	"github.com/ofkm/goobrew/internal/logger"
	"github.com/ofkm/goobrew/internal/ui"
	"github.com/spf13/cobra"
)

// listCmd represents the list command.
// It retrieves and displays all installed Homebrew packages with detailed information
// including version numbers, descriptions, and status indicators (outdated, pinned).
// The list is formatted with colors and icons for better readability.
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List installed packages",
	Long:    `List all installed Homebrew packages with detailed information.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		logger.Log.Info("fetching installed packages")

		formulae, err := client.GetInstalledFormulae(ctx)
		if err != nil {
			ui.PrintError("Failed to get installed packages: " + err.Error())
			logger.Log.Error("failed to get installed packages", "error", err)
			os.Exit(1)
		}

		ui.PrintInstalledList(formulae)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
