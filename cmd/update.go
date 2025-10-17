package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ofkm/goobrew/internal/logger"
	"github.com/ofkm/goobrew/internal/ui"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command.
// It updates Homebrew itself and refreshes the list of available formulae and casks
// from GitHub. This command should be run periodically to ensure access to the
// latest package information.
var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"up"},
	Short:   "Update Homebrew and formulae",
	Long:    `Update Homebrew itself and all formulae from GitHub.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		fmt.Printf("\n%s %sUpdating Homebrew...%s\n\n", ui.IconUpdate, ui.Bold, ui.Reset)
		start := time.Now()

		logger.Log.Info("updating homebrew")

		if err := client.Update(ctx); err != nil {
			elapsed := time.Since(start)
			ui.PrintError(fmt.Sprintf("Update failed (took %s): %v", ui.FormatDuration(elapsed), err))
			logger.Log.Error("update failed", "error", err)
			os.Exit(1)
		}

		elapsed := time.Since(start)
		fmt.Printf("\n%s Update completed in %s%s%s\n\n",
			ui.IconSuccess, ui.Green, ui.FormatDuration(elapsed), ui.Reset)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
