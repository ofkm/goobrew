package cmd

import (
	"context"
	"os"

	"github.com/ofkm/goobrew/internal/logger"
	"github.com/ofkm/goobrew/internal/ui"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command.
// It displays detailed information about a specific package including version,
// dependencies, installation status, and any caveats. The information is
// retrieved from Homebrew's JSON API and formatted for easy reading.
var infoCmd = &cobra.Command{
	Use:   "info [package]",
	Short: "Display package information",
	Long:  `Display detailed information about a package in a beautiful format.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		pkgName := args[0]

		logger.Log.Info("fetching package info", "package", pkgName)

		formula, err := client.GetFormula(ctx, pkgName)
		if err != nil {
			ui.PrintError("Failed to get package info: " + err.Error())
			logger.Log.Error("failed to get formula info", "error", err, "package", pkgName)
			os.Exit(1)
		}

		ui.PrintFormulaInfo(formula)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
