package cmd

import (
	"fmt"

	"github.com/ofkm/goobrew/internal/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command.
// It displays the version information for goobrew including the semantic version,
// git commit hash, and build time.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display version information including git commit and build time.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.GetFullVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
