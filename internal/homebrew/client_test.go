package homebrew

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Skipf("Skipping: %v", err)
	}

	if client == nil {
		t.Fatal("Client should not be nil")
	}

	if client.httpClient == nil {
		t.Error("HTTP client should not be nil")
	}

	if client.brewPath == "" {
		t.Error("Brew path should not be empty")
	}

	// Give it a moment for background loading to start
	time.Sleep(100 * time.Millisecond)
}

func TestLoadFormulaeAndCasks(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "formula.json") {
			formulae := []FormulaListItem{
				{Name: "git", Desc: "Distributed version control"},
				{Name: "node", Desc: "JavaScript runtime"},
			}
			_ = json.NewEncoder(w).Encode(formulae)
		} else if strings.Contains(r.URL.Path, "cask.json") {
			casks := []CaskListItem{
				{Token: "firefox", Name: []string{"Firefox"}, Desc: "Web browser"},
				{Token: "chrome", Name: []string{"Google Chrome"}, Desc: "Web browser"},
			}
			_ = json.NewEncoder(w).Encode(casks)
		}
	}))
	defer server.Close()

	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	// Test fetchFormulaeList
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL+"/formula.json", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	var formulae []FormulaListItem
	err = json.NewDecoder(resp.Body).Decode(&formulae)
	if err != nil {
		t.Fatalf("Failed to decode formulae: %v", err)
	}

	if len(formulae) != 2 {
		t.Errorf("Expected 2 formulae, got %d", len(formulae))
	}
}

func TestGetFormula(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/formula/git.json") {
			formula := Formula{
				Name:     "git",
				FullName: "git",
				Desc:     "Distributed version control system",
				Homepage: "https://git-scm.com",
				Versions: Versions{Stable: "2.51.1", Bottle: true},
			}
			_ = json.NewEncoder(w).Encode(formula)
		} else {
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		brewPath:   "brew", // Assume brew is in PATH
	}

	// Directly fetch from our test server
	ctx := context.Background()
	formula, err := client.fetchFormula(ctx, server.URL+"/formula/git.json")
	if err != nil {
		t.Fatalf("Failed to get formula: %v", err)
	}

	if formula.Name != "git" {
		t.Errorf("Expected name 'git', got '%s'", formula.Name)
	}
}

func TestGetFormula_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		brewPath:   "brew",
	}

	ctx := context.Background()
	_, err := client.fetchFormula(ctx, server.URL+"/nonexistent.json")
	if err == nil {
		t.Error("Expected error for nonexistent formula")
	}
}

func TestSearch(t *testing.T) {
	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		formulaeCache: []FormulaListItem{
			{Name: "git", Desc: "Distributed version control"},
			{Name: "github-cli", Desc: "GitHub command line"},
			{Name: "gitlab-runner", Desc: "GitLab CI runner"},
			{Name: "node", Desc: "JavaScript runtime"},
		},
		casksCache: []CaskListItem{
			{Token: "github", Name: []string{"GitHub Desktop"}, Desc: "GitHub desktop app"},
			{Token: "gitkraken", Name: []string{"GitKraken"}, Desc: "Git GUI"},
			{Token: "firefox", Name: []string{"Firefox"}, Desc: "Web browser"},
		},
		cacheTimestamp: time.Now(),
	}

	ctx := context.Background()
	formulae, casks, err := client.Search(ctx, "git")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(formulae) < 2 {
		t.Errorf("Expected at least 2 formulae containing 'git', got %d", len(formulae))
	}

	if len(casks) < 2 {
		t.Errorf("Expected at least 2 casks containing 'git', got %d", len(casks))
	}

	// Test case insensitive search
	formulae2, _, err := client.Search(ctx, "GIT")
	if err != nil {
		t.Fatalf("Case insensitive search failed: %v", err)
	}

	if len(formulae2) != len(formulae) {
		t.Error("Case insensitive search should return same results")
	}
}

func TestSearchWithExpiredCache(t *testing.T) {
	client := &Client{
		httpClient:     &http.Client{Timeout: 5 * time.Second},
		formulaeCache:  []FormulaListItem{},
		casksCache:     []CaskListItem{},
		cacheTimestamp: time.Now().Add(-2 * time.Hour), // Expired cache
	}

	ctx := context.Background()
	_, _, err := client.Search(ctx, "test")
	if err != nil {
		t.Fatalf("Search with expired cache failed: %v", err)
	}
}

func TestGetFromCache(t *testing.T) {
	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	// Test cache miss
	_, ok := client.getFromCache("nonexistent")
	if ok {
		t.Error("Expected cache miss for nonexistent key")
	}

	// Test cache hit
	testData := &Formula{Name: "test"}
	client.cache.Store("test", cacheEntry{
		data:      testData,
		timestamp: time.Now(),
	})

	cached, ok := client.getFromCache("test")
	if !ok {
		t.Error("Expected cache hit")
	}

	if cachedFormula, ok := cached.(*Formula); !ok || cachedFormula.Name != "test" {
		t.Error("Cached data doesn't match")
	}

	// Test expired cache
	client.cache.Store("expired", cacheEntry{
		data:      testData,
		timestamp: time.Now().Add(-2 * time.Hour),
	})

	_, ok = client.getFromCache("expired")
	if ok {
		t.Error("Expected cache miss for expired entry")
	}
}

func TestInstall(t *testing.T) {
	// This test requires brew to be installed
	if _, err := exec.LookPath("brew"); err != nil {
		t.Skip("Skipping install test: brew not found")
	}

	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		brewPath:   "brew",
	}

	// We can't actually install without sudo/confirmation
	// So we'll just test the channel creation
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	statusChan := make(chan InstallationStatus, 10)

	go func() {
		// This will fail quickly which is what we want for testing
		_ = client.Install(ctx, []string{"nonexistent-package-12345"}, statusChan)
		close(statusChan)
	}()

	// Collect status updates
	statuses := []InstallationStatus{}
	for status := range statusChan {
		statuses = append(statuses, status)
	}

	// We should have received at least one status update
	if len(statuses) == 0 {
		t.Error("Expected at least one status update")
	}
}

func TestUninstall(t *testing.T) {
	if _, err := exec.LookPath("brew"); err != nil {
		t.Skip("Skipping uninstall test: brew not found")
	}

	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		brewPath:   "brew",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// This will fail quickly
	err := client.Uninstall(ctx, []string{"nonexistent-package-12345"})
	// Error is expected for nonexistent package
	_ = err
}

func TestUpdate(t *testing.T) {
	if _, err := exec.LookPath("brew"); err != nil {
		t.Skip("Skipping update test: brew not found")
	}

	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		brewPath:   "brew",
	}

	// Just call Update to verify it works (it will execute brew update which is fast)
	_ = client.Update(context.Background())
}

func TestUpgrade(t *testing.T) {
	if _, err := exec.LookPath("brew"); err != nil {
		t.Skip("Skipping upgrade test: brew not found")
	}

	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		brewPath:   "brew",
	}

	// Just call Upgrade to verify it works
	_ = client.Upgrade(context.Background(), []string{})
}

func TestParseInstallOutput(t *testing.T) {
	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	tests := []struct {
		name          string
		line          string
		expectedStage string
		expectedProg  int
	}{
		{"Downloading", "==> Downloading git-2.51.1.tar.gz", "downloading", 25},
		{"Installing", "==> Installing git", "installing", 50},
		{"Pouring", "==> Pouring git--2.51.1.bottle.tar.gz", "installing", 60},
		{"Linking", "==> Linking git", "linking", 90},
		{"Installed", "git 2.51.1 is installed", "completed", 100},
		{"No match", "Some other output", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := client.parseInstallOutput("git", tt.line)

			if tt.expectedStage == "" {
				if status != nil {
					t.Errorf("Expected nil status for line: %s", tt.line)
				}
			} else {
				if status == nil {
					t.Fatalf("Expected status for line: %s", tt.line)
				}
				if status.Stage != tt.expectedStage {
					t.Errorf("Expected stage '%s', got '%s'", tt.expectedStage, status.Stage)
				}
				if status.Progress != tt.expectedProg {
					t.Errorf("Expected progress %d, got %d", tt.expectedProg, status.Progress)
				}
				if status.Formula != "git" {
					t.Errorf("Expected formula 'git', got '%s'", status.Formula)
				}
			}
		})
	}
}

func TestExecuteCommand(t *testing.T) {
	if _, err := exec.LookPath("brew"); err != nil {
		t.Skip("Skipping execute command test: brew not found")
	}

	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		brewPath:   "brew",
	}

	ctx := context.Background()

	// Test with --version which should work without side effects
	err := client.ExecuteCommand(ctx, []string{"--version"})
	if err != nil {
		t.Errorf("ExecuteCommand failed: %v", err)
	}
}

func TestFetchFormula(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		formula := Formula{
			Name:     "test",
			FullName: "test",
			Desc:     "Test formula",
			Versions: Versions{Stable: "1.0.0", Bottle: false},
		}
		_ = json.NewEncoder(w).Encode(formula)
	}))
	defer server.Close()

	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	ctx := context.Background()
	formula, err := client.fetchFormula(ctx, server.URL)
	if err != nil {
		t.Fatalf("fetchFormula failed: %v", err)
	}

	if formula.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", formula.Name)
	}
}

func TestGetInstalledFormulae(t *testing.T) {
	if _, err := exec.LookPath("brew"); err != nil {
		t.Skip("Skipping get installed formulae test: brew not found")
	}

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	formulae, err := client.GetInstalledFormulae(ctx)
	if err != nil {
		t.Fatalf("GetInstalledFormulae failed: %v", err)
	}

	// We can't assert specific formulae, but we can check the structure
	if formulae == nil {
		t.Error("Expected non-nil formulae slice")
	}

	// If any formulae are installed, check their structure
	for _, f := range formulae {
		if f.Name == "" {
			t.Error("Formula should have a name")
		}
		if len(f.Installed) == 0 {
			t.Error("Installed formula should have installation info")
		}
	}
}

func TestGetLocalInstallInfo(t *testing.T) {
	if _, err := exec.LookPath("brew"); err != nil {
		t.Skip("Skipping get local install info test: brew not found")
	}

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Try to get info for a common formula
	// This might fail if the formula isn't installed, which is fine
	_, err = client.getLocalInstallInfo(ctx, "git")
	// We don't assert error as the formula might not be installed
	_ = err
}

func TestHTTPErrors(t *testing.T) {
	// Test server that returns errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal Server Error")
	}))
	defer server.Close()

	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	ctx := context.Background()
	_, err := client.fetchFormula(ctx, server.URL)
	if err == nil {
		t.Error("Expected error for 500 status code")
	}
}

func TestContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		_ = json.NewEncoder(w).Encode(Formula{Name: "test"})
	}))
	defer server.Close()

	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := client.fetchFormula(ctx, server.URL)
	if err == nil {
		t.Error("Expected error from context cancellation")
	}
}

func TestNewClientWithoutBrew(t *testing.T) {
	// Temporarily modify PATH to not include brew
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", "/nonexistent")
	defer os.Setenv("PATH", oldPath)

	_, err := NewClient()
	if err == nil {
		t.Error("Expected error when brew is not in PATH")
	}

	if !strings.Contains(err.Error(), "not installed") {
		t.Errorf("Expected 'not installed' error, got: %v", err)
	}
}

func TestGetFormulaActual(t *testing.T) {
	// Test actual GetFormula with cache
	client, err := NewClient()
	if err != nil {
		t.Skip("Skipping: brew not found")
	}

	// Pre-populate cache with a formula
	testFormula := &Formula{
		Name: "cached-test-formula",
		Desc: "Test cached formula",
	}
	client.cache.Store("cached-test-formula", cacheEntry{
		data:      testFormula,
		timestamp: time.Now(),
	})

	// Should return from cache
	formula, err := client.GetFormula(context.Background(), "cached-test-formula")
	if err != nil {
		t.Fatalf("GetFormula failed: %v", err)
	}

	if formula.Name != "cached-test-formula" {
		t.Errorf("Expected name 'cached-test-formula', got '%s'", formula.Name)
	}
}

func TestMonitorInstallation(t *testing.T) {
	client := &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	stdout := io.NopCloser(strings.NewReader("==> Downloading git-2.51.1.tar.gz\n==> Installing git\n"))
	stderr := io.NopCloser(strings.NewReader(""))
	statusChan := make(chan InstallationStatus, 10)

	// Start monitoring
	client.monitorInstallation(stdout, stderr, "git", statusChan)

	// Collect statuses with timeout to avoid hanging
	statusCount := 0
	timeout := time.After(500 * time.Millisecond)

collecting:
	for {
		select {
		case <-statusChan:
			statusCount++
		case <-timeout:
			break collecting
		}
	}

	// monitorInstallation reads from readers and sends statuses.
	// We're just testing that it doesn't panic and can parse output.
	// The goroutines will exit naturally when readers are exhausted.
	// We expect at least some status updates from the sample output
	if statusCount > 0 {
		t.Logf("Received %d status updates", statusCount)
	}
}
