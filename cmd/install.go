package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ofkm/goobrew/internal/homebrew"
	"github.com/ofkm/goobrew/internal/logger"
	"github.com/ofkm/goobrew/internal/ui"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:     "install [package...]",
	Aliases: []string{"i"},
	Short:   "Install packages",
	Long:    `Install one or more Homebrew packages with beautiful progress tracking.`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		fmt.Printf("\n%s %sInstalling packages:%s %s\n\n",
			ui.IconBeer, ui.Bold, ui.Reset, strings.Join(args, ", "))

		start := time.Now()
		statusChan := make(chan homebrew.InstallationStatus, 100)

		// Start installation in background
		go func() {
			defer close(statusChan)
			if err := client.Install(ctx, args, statusChan); err != nil {
				logger.Log.Error("installation failed", "error", err)
			}
		}()

		// Monitor progress
		lastPkg := ""
		for status := range statusChan {
			if status.Formula != lastPkg {
				if lastPkg != "" {
					fmt.Println() // New line for new package
				}
				lastPkg = status.Formula
			}

			ui.PrintInstallProgress(status)

			if status.Stage == "completed" {
				fmt.Println() // New line after completion
				ui.PrintSuccess(fmt.Sprintf("%s installed successfully", status.Formula))
			} else if status.Stage == "failed" {
				fmt.Println() // New line after failure
				ui.PrintError(fmt.Sprintf("Failed to install %s: %v", status.Formula, status.Error))
			}
		}

		elapsed := time.Since(start)
		fmt.Printf("\n%s Installation completed in %s%s%s\n\n",
			ui.IconSparkles, ui.Green, ui.FormatDuration(elapsed), ui.Reset)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
