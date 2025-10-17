package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ofkm/goobrew/internal/logger"
	"github.com/ofkm/goobrew/internal/ui"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command.
// It removes one or more packages from the system using Homebrew's uninstall
// functionality. The command accepts multiple package names and displays
// progress and timing information during the uninstallation process.
var uninstallCmd = &cobra.Command{
	Use:     "uninstall [package...]",
	Aliases: []string{"remove", "rm"},
	Short:   "Uninstall packages",
	Long:    `Uninstall one or more Homebrew packages.`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		fmt.Printf("\n%s %sUninstalling packages:%s %s\n\n",
			ui.IconTrash, ui.Bold, ui.Reset, strings.Join(args, ", "))

		start := time.Now()
		logger.Log.Info("uninstalling packages", "packages", args)

		if err := client.Uninstall(ctx, args); err != nil {
			elapsed := time.Since(start)
			ui.PrintError(fmt.Sprintf("Uninstallation failed (took %s): %v", ui.FormatDuration(elapsed), err))
			logger.Log.Error("uninstallation failed", "error", err)
			os.Exit(1)
		}

		elapsed := time.Since(start)
		fmt.Printf("\n%s Uninstallation completed in %s%s%s\n\n",
			ui.IconSuccess, ui.Green, ui.FormatDuration(elapsed), ui.Reset)
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
