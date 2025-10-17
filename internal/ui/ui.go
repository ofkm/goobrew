// Package ui provides user interface utilities for goobrew.
// It includes color constants, Nerd Font icons, and formatting functions for
// displaying information in a beautiful and user-friendly way.
package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/ofkm/goobrew/internal/homebrew"
)

// Colors are ANSI escape codes for terminal text formatting.
const (
	Reset   = "\033[0m"  // Reset resets all formatting
	Bold    = "\033[1m"  // Bold makes text bold
	Red     = "\033[31m" // Red colors text red
	Green   = "\033[32m" // Green colors text green
	Yellow  = "\033[33m" // Yellow colors text yellow
	Blue    = "\033[34m" // Blue colors text blue
	Magenta = "\033[35m" // Magenta colors text magenta
	Cyan    = "\033[36m" // Cyan colors text cyan
	Gray    = "\033[90m" // Gray colors text gray
)

// Symbols are Nerd Font icons used throughout the UI.
// These icons require a Nerd Font patched font to display correctly.
const (
	IconBeer     = "\uf0f8"     // IconBeer represents Homebrew (nf-dev-homebrew)
	IconPackage  = "\U000f0317" // IconPackage represents packages (nf-md-package)
	IconSearch   = "\uf002"     // IconSearch represents search operations (nf-fa-search)
	IconInfo     = "\uf05a"     // IconInfo represents information (nf-fa-info_circle)
	IconSuccess  = "\uf058"     // IconSuccess represents successful operations (nf-fa-check_circle)
	IconError    = "\uf057"     // IconError represents errors (nf-fa-times_circle)
	IconWarning  = "\uf071"     // IconWarning represents warnings (nf-fa-exclamation_triangle)
	IconDownload = "\uf019"     // IconDownload represents downloads (nf-fa-download)
	IconInstall  = "\uf013"     // IconInstall represents installation (nf-fa-cog)
	IconLink     = "\uf0c1"     // IconLink represents linking (nf-fa-link)
	IconUpdate   = "\uf021"     // IconUpdate represents updates (nf-fa-refresh)
	IconTrash    = "\uf1f8"     // IconTrash represents uninstallation (nf-fa-trash)
	IconSparkles = "\uf005"     // IconSparkles represents completion (nf-fa-star)
	IconRocket   = "\uf135"     // IconRocket represents upgrades (nf-fa-rocket)
)

// FormatDuration formats a time.Duration into a human-readable string.
// It returns durations less than a second as "< 1s", seconds as "Xs",
// minutes as "Xm Ys", and hours as "Xh Ym".
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return "< 1s"
	}
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm %ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh %dm", int(d.Hours()), int(d.Minutes())%60)
}

// FormatSize formats a byte count into a human-readable string using
// binary units (KiB, MiB, GiB, etc.). Values less than 1024 bytes are
// displayed as bytes, larger values are automatically scaled to the
// appropriate unit (KB, MB, GB, etc.).
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// PrintFormulaInfo displays detailed information about a Homebrew formula.
// It prints the formula name, description, homepage, version, license,
// installation status, dependencies, build dependencies, and any caveats.
// All output is formatted with colors and icons for readability.
func PrintFormulaInfo(formula *homebrew.Formula) {
	fmt.Printf("\n%s %s%s%s\n", IconInfo, Bold, formula.Name, Reset)

	if formula.Desc != "" {
		fmt.Printf("  %s\n", formula.Desc)
	}

	fmt.Printf("\n  %sHomepage:%s %s\n", Cyan, Reset, formula.Homepage)

	if formula.Versions.Stable != "" {
		fmt.Printf("  %sVersion:%s  %s\n", Cyan, Reset, formula.Versions.Stable)
	}

	if formula.License != "" {
		fmt.Printf("  %sLicense:%s  %s\n", Cyan, Reset, formula.License)
	}

	// Installation status
	if len(formula.Installed) > 0 {
		latest := formula.Installed[len(formula.Installed)-1]
		installTime := time.Unix(latest.Time, 0)
		fmt.Printf("\n  %s%sInstalled:%s %s %s(on %s)%s\n",
			Green, Bold, Reset, latest.Version, Gray, installTime.Format("Jan 02, 2006"), Reset)

		if latest.PouredFromBottle {
			fmt.Printf("  %sInstalled from:%s bottle\n", Cyan, Reset)
		} else {
			fmt.Printf("  %sInstalled from:%s source\n", Cyan, Reset)
		}
	} else {
		fmt.Printf("\n  %sNot installed%s\n", Yellow, Reset)
	}

	// Dependencies
	if len(formula.Dependencies) > 0 {
		fmt.Printf("\n  %sDependencies:%s\n", Cyan, Reset)
		for _, dep := range formula.Dependencies {
			fmt.Printf("    • %s\n", dep)
		}
	}

	// Build dependencies
	if len(formula.BuildDependencies) > 0 {
		fmt.Printf("\n  %sBuild Dependencies:%s\n", Cyan, Reset)
		for _, dep := range formula.BuildDependencies {
			fmt.Printf("    • %s\n", dep)
		}
	}

	// Caveats
	if formula.Caveats != "" {
		fmt.Printf("\n  %s%sℹ️  Caveats:%s\n", Yellow, Bold, Reset)
		caveats := strings.TrimSpace(formula.Caveats)
		for _, line := range strings.Split(caveats, "\n") {
			fmt.Printf("  %s\n", line)
		}
	}

	fmt.Println()
}

// PrintSearchResults displays search results for formulae and casks.
// It separates formulae and casks into distinct sections with appropriate
// icons and colors. If no results are found, it displays a warning message.
// The function also shows the total count of matching packages.
func PrintSearchResults(formulae, casks []string) {
	if len(formulae) > 0 {
		fmt.Printf("\n%s %s%sFormulae%s\n", IconPackage, Bold, Green, Reset)
		for _, f := range formulae {
			fmt.Printf("  • %s\n", f)
		}
	}

	if len(casks) > 0 {
		fmt.Printf("\n%s %s%sCasks%s\n", IconPackage, Bold, Cyan, Reset)
		for _, c := range casks {
			fmt.Printf("  • %s\n", c)
		}
	}

	total := len(formulae) + len(casks)
	if total == 0 {
		fmt.Printf("\n%s No results found\n", IconWarning)
	} else {
		fmt.Printf("\n%sTotal:%s %d results\n", Gray, Reset, total)
	}
	fmt.Println()
}

// PrintInstalledList displays a list of all installed packages.
// Each package is shown with its name, version, and a status indicator
// (green for up-to-date, yellow for outdated, blue for pinned).
// If no packages are installed, it displays a warning message.
func PrintInstalledList(formulae []homebrew.Formula) {
	if len(formulae) == 0 {
		fmt.Printf("\n%s No packages installed\n\n", IconWarning)
		return
	}

	fmt.Printf("\n%s %s%sInstalled Packages%s (%d total)\n\n", IconPackage, Bold, Green, Reset, len(formulae))

	for _, f := range formulae {
		version := ""
		if len(f.Installed) > 0 {
			version = f.Installed[len(f.Installed)-1].Version
		}

		statusIcon := Green + "●" + Reset
		if f.Outdated {
			statusIcon = Yellow + "●" + Reset
		}
		if f.Pinned {
			statusIcon = Blue + "📌" + Reset
		}

		fmt.Printf("  %s %s%-30s%s %s%s%s", statusIcon, Cyan, f.Name, Reset, Gray, version, Reset)

		if f.Desc != "" {
			desc := f.Desc
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}
			fmt.Printf(" - %s", desc)
		}
		fmt.Println()
	}

	fmt.Println()
}

// PrintInstallProgress displays real-time installation progress.
// It shows the package name, current installation stage (downloading,
// installing, linking, completed, or failed), elapsed time, and any
// errors that occurred. The output uses different icons and colors
// based on the installation stage.
func PrintInstallProgress(status homebrew.InstallationStatus) {
	elapsed := time.Since(status.StartTime)

	icon := IconInstall
	color := Cyan

	switch status.Stage {
	case "downloading":
		icon = IconDownload
		color = Blue
	case "installing":
		icon = IconInstall
		color = Yellow
	case "linking":
		icon = IconLink
		color = Magenta
	case "completed":
		icon = IconSuccess
		color = Green
	case "failed":
		icon = IconError
		color = Red
	}

	fmt.Printf("\r%s %s%s%s %s%-20s%s [%s]",
		icon, color, status.Formula, Reset, Gray, status.Stage, Reset, FormatDuration(elapsed))

	if status.Error != nil {
		fmt.Printf(" - %s%s%s", Red, status.Error, Reset)
	}
}

// PrintSuccess displays a success message with a checkmark icon and green color.
func PrintSuccess(message string) {
	fmt.Printf("%s %s%s%s\n", IconSuccess, Green, message, Reset)
}

// PrintError displays an error message with an error icon and red color.
func PrintError(message string) {
	fmt.Printf("%s %s%s%s\n", IconError, Red, message, Reset)
}

// PrintWarning displays a warning message with a warning icon and yellow color.
func PrintWarning(message string) {
	fmt.Printf("%s %s%s%s\n", IconWarning, Yellow, message, Reset)
}

// PrintInfo displays an informational message with an info icon.
func PrintInfo(message string) {
	fmt.Printf("%s %s\n", IconInfo, message)
}

// ProgressBar creates a visual progress bar string.
// It takes the current progress value, total value, and desired width in characters.
// Returns a colored progress bar with percentage. If total is 0, returns a full bar.
// The filled portion is displayed in green and the empty portion in gray.
func ProgressBar(current, total int, width int) string {
	if total == 0 {
		return strings.Repeat("━", width)
	}

	percent := float64(current) / float64(total)
	filled := int(percent * float64(width))

	bar := strings.Repeat("━", filled)
	empty := strings.Repeat("╌", width-filled)

	return fmt.Sprintf("%s%s%s%s %3.0f%%", Green, bar, Gray, empty, percent*100)
}
