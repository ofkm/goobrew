package ui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ofkm/goobrew/internal/homebrew"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		contains string
	}{
		{"Less than 1 second", 500 * time.Millisecond, "< 1s"},
		{"Exactly 1 second", 1 * time.Second, "1s"},
		{"30 seconds", 30 * time.Second, "30s"},
		{"1 minute 30 seconds", 90 * time.Second, "1m"},
		{"1 hour 30 minutes", 90 * time.Minute, "1h"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.duration)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("FormatDuration(%v) = %s, expected to contain '%s'", tt.duration, result, tt.contains)
			}
		})
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{"Less than KB", 500, "500 B"},
		{"1 KB", 1024, "1.0 KB"},
		{"1 MB", 1024 * 1024, "1.0 MB"},
		{"1.5 MB", 1536 * 1024, "1.5 MB"},
		{"1 GB", 1024 * 1024 * 1024, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatSize(tt.bytes)
			if result != tt.expected {
				t.Errorf("FormatSize(%d) = %s, expected %s", tt.bytes, result, tt.expected)
			}
		})
	}
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestPrintFormulaInfo(t *testing.T) {
	formula := &homebrew.Formula{
		Name:     "git",
		Desc:     "Distributed revision control system",
		Homepage: "https://git-scm.com",
		License:  "GPL-2.0-only",
		Versions: homebrew.Versions{
			Stable: "2.51.1",
			Bottle: true,
		},
		Dependencies:      []string{"gettext", "pcre2"},
		BuildDependencies: []string{"go"},
		Installed: []homebrew.InstalledInfo{
			{
				Version:          "2.51.1",
				PouredFromBottle: true,
				Time:             time.Now().Unix(),
			},
		},
		Caveats: "Some important caveats",
	}

	output := captureOutput(func() {
		PrintFormulaInfo(formula)
	})

	if !strings.Contains(output, "git") {
		t.Error("Output should contain formula name")
	}
	if !strings.Contains(output, "Distributed revision control system") {
		t.Error("Output should contain description")
	}
	if !strings.Contains(output, "2.51.1") {
		t.Error("Output should contain version")
	}
	if !strings.Contains(output, "gettext") {
		t.Error("Output should contain dependencies")
	}
}

func TestPrintSearchResults(t *testing.T) {
	formulae := []string{"git", "github-cli", "gitlab-runner"}
	casks := []string{"github", "gitkraken"}

	output := captureOutput(func() {
		PrintSearchResults(formulae, casks)
	})

	if !strings.Contains(output, "git") {
		t.Error("Output should contain formulae")
	}
	if !strings.Contains(output, "github") {
		t.Error("Output should contain casks")
	}
	if !strings.Contains(output, "5 results") {
		t.Error("Output should show total count")
	}
}

func TestPrintSearchResults_Empty(t *testing.T) {
	output := captureOutput(func() {
		PrintSearchResults([]string{}, []string{})
	})

	if !strings.Contains(output, "No results") {
		t.Error("Output should indicate no results found")
	}
}

func TestPrintInstalledList(t *testing.T) {
	formulae := []homebrew.Formula{
		{
			Name: "git",
			Desc: "Distributed revision control system",
			Installed: []homebrew.InstalledInfo{
				{Version: "2.51.1"},
			},
			Outdated: false,
			Pinned:   false,
		},
		{
			Name: "node",
			Desc: "JavaScript runtime",
			Installed: []homebrew.InstalledInfo{
				{Version: "20.0.0"},
			},
			Outdated: true,
			Pinned:   false,
		},
		{
			Name: "python",
			Desc: "Interpreted, interactive, object-oriented programming language",
			Installed: []homebrew.InstalledInfo{
				{Version: "3.12.0"},
			},
			Outdated: false,
			Pinned:   true,
		},
	}

	output := captureOutput(func() {
		PrintInstalledList(formulae)
	})

	if !strings.Contains(output, "git") {
		t.Error("Output should contain git")
	}
	if !strings.Contains(output, "node") {
		t.Error("Output should contain node")
	}
	if !strings.Contains(output, "python") {
		t.Error("Output should contain python")
	}
	if !strings.Contains(output, "3 total") {
		t.Error("Output should show total count")
	}
}

func TestPrintInstalledList_Empty(t *testing.T) {
	output := captureOutput(func() {
		PrintInstalledList([]homebrew.Formula{})
	})

	if !strings.Contains(output, "No packages") {
		t.Error("Output should indicate no packages installed")
	}
}

func TestPrintInstallProgress(t *testing.T) {
	tests := []struct {
		name   string
		status homebrew.InstallationStatus
	}{
		{
			name: "Downloading",
			status: homebrew.InstallationStatus{
				Formula:   "git",
				Stage:     "downloading",
				Progress:  25,
				StartTime: time.Now(),
			},
		},
		{
			name: "Installing",
			status: homebrew.InstallationStatus{
				Formula:   "git",
				Stage:     "installing",
				Progress:  50,
				StartTime: time.Now(),
			},
		},
		{
			name: "Completed",
			status: homebrew.InstallationStatus{
				Formula:   "git",
				Stage:     "completed",
				Progress:  100,
				StartTime: time.Now().Add(-5 * time.Second),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just ensure it doesn't panic
			PrintInstallProgress(tt.status)
		})
	}
}

func TestPrintMessages(t *testing.T) {
	tests := []struct {
		name string
		fn   func(string)
		msg  string
	}{
		{"Success", PrintSuccess, "Installation complete"},
		{"Error", PrintError, "Failed to install"},
		{"Warning", PrintWarning, "Package outdated"},
		{"Info", PrintInfo, "Additional information"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just ensure they don't panic
			tt.fn(tt.msg)
		})
	}
}

func TestProgressBar(t *testing.T) {
	tests := []struct {
		name     string
		current  int
		total    int
		width    int
		contains string
	}{
		{"50%", 50, 100, 20, "50%"},
		{"100%", 100, 100, 20, "100%"},
		{"0%", 0, 100, 20, "0%"},
		{"Zero total", 0, 0, 20, "━━━━━━━━━━━━━━━━━━━━"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProgressBar(tt.current, tt.total, tt.width)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("ProgressBar(%d, %d, %d) = %s, expected to contain '%s'",
					tt.current, tt.total, tt.width, result, tt.contains)
			}
		})
	}
}

func TestAllColorConstants(t *testing.T) {
	// Ensure all color constants are defined
	colors := []string{Reset, Bold, Red, Green, Yellow, Blue, Magenta, Cyan, Gray}
	for i, color := range colors {
		if color == "" {
			t.Errorf("Color constant at index %d is empty", i)
		}
	}
}

func TestAllIconConstants(t *testing.T) {
	// Ensure all icon constants are defined
	icons := []string{
		IconBeer, IconPackage, IconSearch, IconInfo, IconSuccess,
		IconError, IconWarning, IconDownload, IconInstall, IconLink,
		IconUpdate, IconTrash, IconSparkles, IconRocket,
	}
	for i, icon := range icons {
		if icon == "" {
			t.Errorf("Icon constant at index %d is empty", i)
		}
	}
}

func TestPrintFormulaInfo_MinimalFormula(t *testing.T) {
	// Test with minimal formula data
	formula := &homebrew.Formula{
		Name:     "minimal",
		Homepage: "https://example.com",
		Versions: homebrew.Versions{},
	}

	output := captureOutput(func() {
		PrintFormulaInfo(formula)
	})

	if !strings.Contains(output, "minimal") {
		t.Error("Output should contain formula name")
	}
}

func TestPrintFormulaInfo_WithAllFields(t *testing.T) {
	formula := &homebrew.Formula{
		Name:              "full",
		Desc:              "Full test formula",
		Homepage:          "https://example.com",
		License:           "MIT",
		Versions:          homebrew.Versions{Stable: "1.0.0"},
		Dependencies:      []string{"dep1", "dep2"},
		BuildDependencies: []string{"build-dep"},
		Installed: []homebrew.InstalledInfo{
			{
				Version:          "1.0.0",
				PouredFromBottle: false,
				Time:             time.Now().Unix(),
			},
		},
		Caveats: "Line 1\nLine 2\nLine 3",
	}

	output := captureOutput(func() {
		PrintFormulaInfo(formula)
	})

	if !strings.Contains(output, "full") {
		t.Error("Output should contain formula name")
	}
	if !strings.Contains(output, "MIT") {
		t.Error("Output should contain license")
	}
	if !strings.Contains(output, "Line 1") {
		t.Error("Output should contain caveats")
	}
	if !strings.Contains(output, "source") {
		t.Error("Output should indicate installed from source")
	}
}
