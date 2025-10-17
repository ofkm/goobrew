package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// Helper to capture command output
func executeCommand(args ...string) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Reset command for clean test
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}

func TestRootCommand(t *testing.T) {
	output, err := executeCommand("--help")
	if err != nil {
		t.Fatalf("Root command help failed: %v", err)
	}

	if !strings.Contains(output, "goobrew") {
		t.Error("Root command help should mention goobrew")
	}

	if !strings.Contains(output, "Available Commands") {
		t.Error("Root command help should show available commands")
	}
}

func TestSearchCommandHelp(t *testing.T) {
	output, err := executeCommand("search", "--help")
	if err != nil {
		t.Fatalf("Search help failed: %v", err)
	}

	if !strings.Contains(output, "search") {
		t.Error("Search help should mention search")
	}
}

func TestListCommandHelp(t *testing.T) {
	output, err := executeCommand("list", "--help")
	if err != nil {
		t.Fatalf("List help failed: %v", err)
	}

	if !strings.Contains(output, "list") {
		t.Error("List help should mention list")
	}
}

func TestInfoCommandHelp(t *testing.T) {
	output, err := executeCommand("info", "--help")
	if err != nil {
		t.Fatalf("Info help failed: %v", err)
	}

	if !strings.Contains(output, "info") {
		t.Error("Info help should mention info")
	}
}

func TestInstallCommandHelp(t *testing.T) {
	output, err := executeCommand("install", "--help")
	if err != nil {
		t.Fatalf("Install help failed: %v", err)
	}

	if !strings.Contains(output, "install") {
		t.Error("Install help should mention install")
	}
}

func TestUninstallCommandHelp(t *testing.T) {
	output, err := executeCommand("uninstall", "--help")
	if err != nil {
		t.Fatalf("Uninstall help failed: %v", err)
	}

	if !strings.Contains(output, "uninstall") {
		t.Error("Uninstall help should mention uninstall")
	}
}

func TestUpdateCommandHelp(t *testing.T) {
	output, err := executeCommand("update", "--help")
	if err != nil {
		t.Fatalf("Update help failed: %v", err)
	}

	if !strings.Contains(output, "update") {
		t.Error("Update help should mention update")
	}
}

func TestUpgradeCommandHelp(t *testing.T) {
	output, err := executeCommand("upgrade", "--help")
	if err != nil {
		t.Fatalf("Upgrade help failed: %v", err)
	}

	if !strings.Contains(output, "upgrade") {
		t.Error("Upgrade help should mention upgrade")
	}
}

func TestVerboseFlag(t *testing.T) {
	output, err := executeCommand("--verbose", "--help")
	if err != nil {
		t.Fatalf("Verbose flag failed: %v", err)
	}

	if !strings.Contains(output, "goobrew") {
		t.Error("Verbose flag should work with help")
	}
}

func TestDebugFlag(t *testing.T) {
	output, err := executeCommand("--debug", "--help")
	if err != nil {
		t.Fatalf("Debug flag failed: %v", err)
	}

	if !strings.Contains(output, "goobrew") {
		t.Error("Debug flag should work with help")
	}
}

func TestSearchCommand_NoArgs(t *testing.T) {
	output, _ := executeCommand("search")

	// Search shows help when no args provided
	if !strings.Contains(output, "search") {
		t.Error("Search with no args should show help")
	}
}

func TestInfoCommand_NoArgs(t *testing.T) {
	output, _ := executeCommand("info")

	// Info shows help when no args provided
	if !strings.Contains(output, "info") {
		t.Error("Info with no args should show help")
	}
}

func TestInstallCommand_NoArgs(t *testing.T) {
	output, _ := executeCommand("install")

	// Install shows help when no args provided
	if !strings.Contains(output, "install") {
		t.Error("Install with no args should show help")
	}
}

func TestUninstallCommand_NoArgs(t *testing.T) {
	output, _ := executeCommand("uninstall")

	// Uninstall shows help when no args provided
	if !strings.Contains(output, "uninstall") {
		t.Error("Uninstall with no args should show help")
	}
}

func TestCommandsExist(t *testing.T) {
	// Ensure all commands are registered
	commands := []string{"search", "list", "info", "install", "uninstall", "update", "upgrade"}
	for _, cmdName := range commands {
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Name() == cmdName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Command '%s' not registered", cmdName)
		}
	}
}
