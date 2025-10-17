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

var upgradeCmd = &cobra.Command{
	Use:   "upgrade [package...]",
	Short: "Upgrade packages",
	Long:  `Upgrade installed packages to their latest versions.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		if len(args) == 0 {
			fmt.Printf("\n%s %sUpgrading all packages...%s\n\n", ui.IconRocket, ui.Bold, ui.Reset)
		} else {
			fmt.Printf("\n%s %sUpgrading packages:%s %v\n\n", ui.IconRocket, ui.Bold, ui.Reset, args)
		}

		start := time.Now()
		logger.Log.Info("upgrading packages", "packages", args)

		if err := client.Upgrade(ctx, args); err != nil {
			elapsed := time.Since(start)
			ui.PrintError(fmt.Sprintf("Upgrade failed (took %s): %v", ui.FormatDuration(elapsed), err))
			logger.Log.Error("upgrade failed", "error", err)
			os.Exit(1)
		}

		elapsed := time.Since(start)
		fmt.Printf("\n%s Upgrade completed in %s%s%s\n\n",
			ui.IconSuccess, ui.Green, ui.FormatDuration(elapsed), ui.Reset)
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
