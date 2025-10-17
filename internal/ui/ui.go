package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/ofkm/goobrew/internal/homebrew"
)

// Colors
const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[90m"
)

// Symbols
const (
	IconBeer     = "🍺"
	IconPackage  = "📦"
	IconSearch   = "🔍"
	IconInfo     = "ℹ️"
	IconSuccess  = "✅"
	IconError    = "❌"
	IconWarning  = "⚠️"
	IconDownload = "⬇️"
	IconInstall  = "⚙️"
	IconLink     = "🔗"
	IconUpdate   = "🔄"
	IconTrash    = "🗑️"
	IconSparkles = "✨"
	IconRocket   = "🚀"
)

// FormatDuration formats a duration in a human-readable way
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

// FormatSize formats bytes in a human-readable way
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

// PrintFormulaInfo displays detailed formula information
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

// PrintSearchResults displays search results
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

// PrintInstalledList displays installed packages
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

// PrintInstallProgress displays installation progress
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

// PrintSuccess displays a success message
func PrintSuccess(message string) {
	fmt.Printf("%s %s%s%s\n", IconSuccess, Green, message, Reset)
}

// PrintError displays an error message
func PrintError(message string) {
	fmt.Printf("%s %s%s%s\n", IconError, Red, message, Reset)
}

// PrintWarning displays a warning message
func PrintWarning(message string) {
	fmt.Printf("%s %s%s%s\n", IconWarning, Yellow, message, Reset)
}

// PrintInfo displays an info message
func PrintInfo(message string) {
	fmt.Printf("%s %s\n", IconInfo, message)
}

// ProgressBar creates a simple progress bar
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
