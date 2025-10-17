package homebrew

import (
	"encoding/json"
	"testing"
	"time"
)

func TestFormulaUnmarshal(t *testing.T) {
	jsonData := `{
		"name": "git",
		"full_name": "git",
		"tap": "homebrew/core",
		"desc": "Distributed revision control system",
		"license": "GPL-2.0-only",
		"homepage": "https://git-scm.com",
		"versions": {
			"stable": "2.51.1",
			"bottle": true
		},
		"revision": 0,
		"version_scheme": 0,
		"bottle": {
			"rebuild": 0,
			"root_url": "https://example.com",
			"files": {}
		},
		"keg_only": false,
		"options": [],
		"build_dependencies": [],
		"dependencies": ["gettext", "pcre2"],
		"test_dependencies": [],
		"recommended_dependencies": [],
		"optional_dependencies": [],
		"uses_from_macos": ["curl", {"zlib": "build"}],
		"uses_from_macos_bounds": [{}, {"since": "monterey"}],
		"requirements": [],
		"conflicts_with": [],
		"conflicts_with_reasons": [],
		"link_overwrite": [],
		"installed": [],
		"pinned": false,
		"outdated": false,
		"deprecated": false,
		"disabled": false,
		"post_install_defined": false,
		"aliases": [],
		"versioned_formulae": []
	}`

	var formula Formula
	err := json.Unmarshal([]byte(jsonData), &formula)
	if err != nil {
		t.Fatalf("Failed to unmarshal formula: %v", err)
	}

	if formula.Name != "git" {
		t.Errorf("Expected name 'git', got '%s'", formula.Name)
	}

	if formula.Desc != "Distributed revision control system" {
		t.Errorf("Expected desc 'Distributed revision control system', got '%s'", formula.Desc)
	}

	if formula.Versions.Stable != "2.51.1" {
		t.Errorf("Expected version '2.51.1', got '%s'", formula.Versions.Stable)
	}

	if len(formula.Dependencies) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(formula.Dependencies))
	}
}

func TestServiceGetKeepAliveBool(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected bool
	}{
		{
			name:     "KeepAlive as true boolean",
			jsonData: `{"name":"test","run_type":"immediate","keep_alive":true}`,
			expected: true,
		},
		{
			name:     "KeepAlive as false boolean",
			jsonData: `{"name":"test","run_type":"immediate","keep_alive":false}`,
			expected: false,
		},
		{
			name:     "KeepAlive as object",
			jsonData: `{"name":"test","run_type":"immediate","keep_alive":{"always":true}}`,
			expected: true,
		},
		{
			name:     "KeepAlive missing",
			jsonData: `{"name":"test","run_type":"immediate"}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var service Service
			err := json.Unmarshal([]byte(tt.jsonData), &service)
			if err != nil {
				t.Fatalf("Failed to unmarshal service: %v", err)
			}

			result := service.GetKeepAliveBool()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestInstallationStatus(t *testing.T) {
	status := InstallationStatus{
		Formula:   "git",
		Stage:     "installing",
		Progress:  50,
		StartTime: time.Now(),
		Error:     nil,
	}

	if status.Formula != "git" {
		t.Errorf("Expected formula 'git', got '%s'", status.Formula)
	}

	if status.Stage != "installing" {
		t.Errorf("Expected stage 'installing', got '%s'", status.Stage)
	}

	if status.Progress != 50 {
		t.Errorf("Expected progress 50, got %d", status.Progress)
	}
}

func TestDependency(t *testing.T) {
	jsonData := `{
		"full_name": "openssl@3",
		"version": "3.0.0",
		"revision": 1,
		"pkg_version": "3.0.0_1",
		"declared_directly": true
	}`

	var dep Dependency
	err := json.Unmarshal([]byte(jsonData), &dep)
	if err != nil {
		t.Fatalf("Failed to unmarshal dependency: %v", err)
	}

	if dep.FullName != "openssl@3" {
		t.Errorf("Expected full_name 'openssl@3', got '%s'", dep.FullName)
	}

	if dep.Version != "3.0.0" {
		t.Errorf("Expected version '3.0.0', got '%s'", dep.Version)
	}
}

func TestInstalledInfo(t *testing.T) {
	jsonData := `{
		"version": "2.51.1",
		"used_options": [],
		"built_as_bottle": true,
		"poured_from_bottle": true,
		"time": 1729123456,
		"runtime_dependencies": [],
		"installed_as_dependency": false,
		"installed_on_request": true
	}`

	var installed InstalledInfo
	err := json.Unmarshal([]byte(jsonData), &installed)
	if err != nil {
		t.Fatalf("Failed to unmarshal installed info: %v", err)
	}

	if installed.Version != "2.51.1" {
		t.Errorf("Expected version '2.51.1', got '%s'", installed.Version)
	}

	if !installed.PouredFromBottle {
		t.Error("Expected PouredFromBottle to be true")
	}

	if !installed.InstalledOnRequest {
		t.Error("Expected InstalledOnRequest to be true")
	}
}

func TestVersions(t *testing.T) {
	jsonData := `{
		"stable": "1.2.3",
		"head": "HEAD",
		"bottle": true
	}`

	var versions Versions
	err := json.Unmarshal([]byte(jsonData), &versions)
	if err != nil {
		t.Fatalf("Failed to unmarshal versions: %v", err)
	}

	if versions.Stable != "1.2.3" {
		t.Errorf("Expected stable '1.2.3', got '%s'", versions.Stable)
	}

	if versions.Head != "HEAD" {
		t.Errorf("Expected head 'HEAD', got '%s'", versions.Head)
	}

	if !versions.Bottle {
		t.Error("Expected Bottle to be true")
	}
}

func TestComplexUsesFromMacos(t *testing.T) {
	jsonData := `{
		"name": "test",
		"full_name": "test",
		"desc": "Test formula",
		"homepage": "https://test.com",
		"versions": {"stable": "1.0.0", "bottle": false},
		"revision": 0,
		"version_scheme": 0,
		"bottle": {"rebuild": 0, "root_url": "", "files": {}},
		"uses_from_macos": [
			"zlib",
			{"flex": "build"},
			{"m4": "build"}
		],
		"uses_from_macos_bounds": [
			{},
			{"since": "monterey"},
			{}
		],
		"keg_only": false,
		"options": [],
		"build_dependencies": [],
		"dependencies": [],
		"test_dependencies": [],
		"recommended_dependencies": [],
		"optional_dependencies": [],
		"requirements": [],
		"conflicts_with": [],
		"conflicts_with_reasons": [],
		"link_overwrite": [],
		"installed": [],
		"pinned": false,
		"outdated": false,
		"deprecated": false,
		"disabled": false,
		"post_install_defined": false,
		"aliases": [],
		"versioned_formulae": []
	}`

	var formula Formula
	err := json.Unmarshal([]byte(jsonData), &formula)
	if err != nil {
		t.Fatalf("Failed to unmarshal formula with complex uses_from_macos: %v", err)
	}

	if len(formula.UsesFromMacos) != 3 {
		t.Errorf("Expected 3 uses_from_macos items, got %d", len(formula.UsesFromMacos))
	}

	if len(formula.UsesFromMacosBounds) != 3 {
		t.Errorf("Expected 3 uses_from_macos_bounds items, got %d", len(formula.UsesFromMacosBounds))
	}
}
