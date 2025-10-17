package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ofkm/goobrew/internal/logger"
	"github.com/ofkm/goobrew/internal/ui"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command.
// It searches for packages (both formulae and casks) in Homebrew repositories
// by matching the search term against package names and descriptions.
// The search is case-insensitive and displays results with beautiful formatting.
var searchCmd = &cobra.Command{
	Use:     "search [term]",
	Aliases: []string{"s"},
	Short:   "Search for packages",
	Long:    `Search for packages in Homebrew repositories with beautiful formatting.`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		searchTerm := args[0]

		fmt.Printf("\n%s %sSearching for:%s %s\n", ui.IconSearch, ui.Bold, ui.Reset, searchTerm)

		logger.Log.Info("searching packages", "term", searchTerm)

		formulae, casks, err := client.Search(ctx, searchTerm)
		if err != nil {
			ui.PrintError("Search failed: " + err.Error())
			logger.Log.Error("search failed", "error", err, "term", searchTerm)
			os.Exit(1)
		}

		ui.PrintSearchResults(formulae, casks)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
