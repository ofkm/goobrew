package homebrew

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/ofkm/goobrew/internal/logger"
)

const (
	HomebrewAPIBase     = "https://formulae.brew.sh/api"
	HomebrewAPIFormulae = "https://formulae.brew.sh/api/formula.json"
	HomebrewAPICasks    = "https://formulae.brew.sh/api/cask.json"
	cacheExpiry         = 1 * time.Hour
	installedCacheKey   = "_installed_formulae"
)

type Client struct {
	httpClient     *http.Client
	cache          sync.Map
	brewPath       string
	formulaeCache  []FormulaListItem // Cache of all formulae names and descriptions
	casksCache     []CaskListItem    // Cache of all cask names and descriptions
	cacheMutex     sync.RWMutex
	cacheTimestamp time.Time
}

// FormulaListItem represents a minimal formula entry for listing/searching
type FormulaListItem struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

// CaskListItem represents a minimal cask entry for listing/searching
type CaskListItem struct {
	Token string   `json:"token"`
	Name  []string `json:"name"`
	Desc  string   `json:"desc"`
}

type cacheEntry struct {
	data      interface{}
	timestamp time.Time
}

// NewClient creates a new Homebrew client
func NewClient() (*Client, error) {
	brewPath, err := exec.LookPath("brew")
	if err != nil {
		return nil, fmt.Errorf("homebrew is not installed. Please install it from https://brew.sh")
	}

	client := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		brewPath: brewPath,
	}

	// Pre-load formulae and casks list in background for faster searches
	go client.loadFormulaeAndCasks(context.Background())

	return client, nil
}

// loadFormulaeAndCasks loads the complete list of formulae and casks from the API in parallel
func (c *Client) loadFormulaeAndCasks(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	logger.Log.Debug("loading formulae and casks from API in parallel")

	// Use channels to load formulae and casks concurrently
	formulaeChan := make(chan []FormulaListItem)
	casksChan := make(chan []CaskListItem)

	// Load formulae in parallel
	go func() {
		if formulae, err := c.fetchFormulaeList(ctx); err == nil {
			logger.Log.Debug("loaded formulae from API", "count", len(formulae))
			formulaeChan <- formulae
		} else {
			logger.Log.Warn("failed to load formulae from API", "error", err)
			formulaeChan <- nil
		}
	}()

	// Load casks in parallel
	go func() {
		if casks, err := c.fetchCasksList(ctx); err == nil {
			logger.Log.Debug("loaded casks from API", "count", len(casks))
			casksChan <- casks
		} else {
			logger.Log.Warn("failed to load casks from API", "error", err)
			casksChan <- nil
		}
	}()

	// Wait for both to complete and update cache atomically
	formulae := <-formulaeChan
	casks := <-casksChan

	c.cacheMutex.Lock()
	if formulae != nil {
		c.formulaeCache = formulae
		c.cacheTimestamp = time.Now()
	}
	if casks != nil {
		c.casksCache = casks
	}
	c.cacheMutex.Unlock()
}

// fetchFormulaeList fetches the complete list of formulae from the API
func (c *Client) fetchFormulaeList(ctx context.Context) ([]FormulaListItem, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, HomebrewAPIFormulae, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var formulae []FormulaListItem
	if err := json.Unmarshal(body, &formulae); err != nil {
		return nil, err
	}

	return formulae, nil
}

// fetchCasksList fetches the complete list of casks from the API
func (c *Client) fetchCasksList(ctx context.Context) ([]CaskListItem, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, HomebrewAPICasks, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var casks []CaskListItem
	if err := json.Unmarshal(body, &casks); err != nil {
		return nil, err
	}

	return casks, nil
}

// GetFormula retrieves information about a specific formula from the web API
func (c *Client) GetFormula(ctx context.Context, name string) (*Formula, error) {
	// Check cache first
	if cached, ok := c.getFromCache(name); ok {
		if formula, ok := cached.(*Formula); ok {
			logger.Log.Debug("using cached formula data", "formula", name)
			return formula, nil
		}
	}

	// Fetch from web API
	url := fmt.Sprintf("%s/formula/%s.json", HomebrewAPIBase, name)
	logger.Log.Debug("fetching formula from web API", "url", url)

	formula, err := c.fetchFormula(ctx, url)
	if err != nil {
		// Try as a cask if formula fetch failed
		caskURL := fmt.Sprintf("%s/cask/%s.json", HomebrewAPIBase, name)
		logger.Log.Debug("trying as cask", "url", caskURL)

		if caskFormula, caskErr := c.fetchFormula(ctx, caskURL); caskErr == nil {
			c.cache.Store(name, cacheEntry{data: caskFormula, timestamp: time.Now()})
			return caskFormula, nil
		}

		return nil, fmt.Errorf("package not found: %s (tried both formula and cask)", name)
	}

	// Merge with local installation info
	if installed, err := c.getLocalInstallInfo(ctx, name); err == nil && len(installed) > 0 {
		formula.Installed = installed
	}

	c.cache.Store(name, cacheEntry{data: formula, timestamp: time.Now()})
	return formula, nil
}

// getLocalInstallInfo gets installation info for a formula from local brew
func (c *Client) getLocalInstallInfo(ctx context.Context, name string) ([]InstalledInfo, error) {
	//nolint:gosec // brewPath is validated at client creation
	cmd := exec.CommandContext(ctx, c.brewPath, "info", "--json=v1", name)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var formulae []Formula
	if err := json.Unmarshal(output, &formulae); err != nil {
		return nil, err
	}

	if len(formulae) > 0 {
		return formulae[0].Installed, nil
	}

	return nil, nil
}

// GetInstalledFormulae retrieves all installed formulae
func (c *Client) GetInstalledFormulae(ctx context.Context) ([]Formula, error) {
	//nolint:gosec // brewPath is validated at client creation
	cmd := exec.CommandContext(ctx, c.brewPath, "info", "--json=v1", "--installed")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get installed formulae: %w", err)
	}

	var formulae []Formula
	if err := json.Unmarshal(output, &formulae); err != nil {
		return nil, fmt.Errorf("failed to parse installed formulae: %w", err)
	}

	return formulae, nil
}

// Search searches for formulae and casks using cached API data with parallel processing
func (c *Client) Search(ctx context.Context, term string) ([]string, []string, error) {
	c.cacheMutex.RLock()
	formulaeCache := c.formulaeCache
	casksCache := c.casksCache
	cacheAge := time.Since(c.cacheTimestamp)
	c.cacheMutex.RUnlock()

	// Reload cache if it's too old or empty
	if cacheAge > cacheExpiry || len(formulaeCache) == 0 || len(casksCache) == 0 {
		logger.Log.Debug("cache expired or empty, reloading")
		c.loadFormulaeAndCasks(ctx)

		// Wait a bit for the reload to complete
		time.Sleep(100 * time.Millisecond)

		c.cacheMutex.RLock()
		formulaeCache = c.formulaeCache
		casksCache = c.casksCache
		c.cacheMutex.RUnlock()
	}

	lowerTerm := strings.ToLower(term)

	// Use channels for concurrent search results
	formulaeChan := make(chan []string)
	casksChan := make(chan []string)

	// Search formulae in parallel
	go func() {
		var results []string
		for _, f := range formulaeCache {
			if strings.Contains(strings.ToLower(f.Name), lowerTerm) ||
				strings.Contains(strings.ToLower(f.Desc), lowerTerm) {
				results = append(results, f.Name)
			}
		}
		formulaeChan <- results
	}()

	// Search casks in parallel
	go func() {
		var results []string
		for _, c := range casksCache {
			if strings.Contains(strings.ToLower(c.Token), lowerTerm) ||
				strings.Contains(strings.ToLower(c.Desc), lowerTerm) {
				results = append(results, c.Token)
			} else {
				// Also search in cask names
				for _, name := range c.Name {
					if strings.Contains(strings.ToLower(name), lowerTerm) {
						results = append(results, c.Token)
						break
					}
				}
			}
		}
		casksChan <- results
	}()

	// Wait for both searches to complete
	formulaeResults := <-formulaeChan
	casksResults := <-casksChan

	logger.Log.Debug("search completed",
		"term", term,
		"formulae_results", len(formulaeResults),
		"casks_results", len(casksResults))

	return formulaeResults, casksResults, nil
}

// Install installs one or more packages using brew command
func (c *Client) Install(ctx context.Context, packages []string, statusChan chan<- InstallationStatus) error {
	for _, pkg := range packages {
		status := InstallationStatus{
			Formula:   pkg,
			Stage:     "starting",
			StartTime: time.Now(),
		}
		statusChan <- status

		//nolint:gosec // brewPath is validated at client creation
		cmd := exec.CommandContext(ctx, c.brewPath, "install", pkg)

		// Create pipes for stdout and stderr
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			status.Stage = "failed"
			status.Error = err
			statusChan <- status
			continue
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			status.Stage = "failed"
			status.Error = err
			statusChan <- status
			continue
		}

		if err := cmd.Start(); err != nil {
			status.Stage = "failed"
			status.Error = err
			statusChan <- status
			continue
		}

		// Monitor output
		go c.monitorInstallation(stdout, stderr, pkg, statusChan)

		if err := cmd.Wait(); err != nil {
			status.Stage = "failed"
			status.Error = err
			statusChan <- status
			continue
		}

		status.Stage = "completed"
		status.Progress = 100
		statusChan <- status
	}

	return nil
}

// Uninstall removes one or more packages
func (c *Client) Uninstall(ctx context.Context, packages []string) error {
	args := append([]string{"uninstall"}, packages...)
	//nolint:gosec // brewPath is validated at client creation, args are package names
	cmd := exec.CommandContext(ctx, c.brewPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Update updates Homebrew
func (c *Client) Update(ctx context.Context) error {
	//nolint:gosec // brewPath is validated at client creation
	cmd := exec.CommandContext(ctx, c.brewPath, "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Upgrade upgrades packages
func (c *Client) Upgrade(ctx context.Context, packages []string) error {
	args := append([]string{"upgrade"}, packages...)
	//nolint:gosec // brewPath is validated at client creation, args are package names
	cmd := exec.CommandContext(ctx, c.brewPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Helper methods

func (c *Client) fetchFormula(ctx context.Context, url string) (*Formula, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var formula Formula
	if err := json.Unmarshal(body, &formula); err != nil {
		return nil, err
	}

	return &formula, nil
}

func (c *Client) getFromCache(key string) (interface{}, bool) {
	if val, ok := c.cache.Load(key); ok {
		entry := val.(cacheEntry)
		if time.Since(entry.timestamp) < cacheExpiry {
			return entry.data, true
		}
		c.cache.Delete(key)
	}
	return nil, false
}

func (c *Client) monitorInstallation(stdout, stderr io.Reader, pkg string, statusChan chan<- InstallationStatus) {
	scanner := func(r io.Reader) {
		buf := make([]byte, 1024)
		for {
			n, err := r.Read(buf)
			if n > 0 {
				line := string(buf[:n])
				status := c.parseInstallOutput(pkg, line)
				if status != nil {
					statusChan <- *status
				}
			}
			if err != nil {
				break
			}
		}
	}

	go scanner(stdout)
	go scanner(stderr)
}

func (c *Client) parseInstallOutput(pkg, line string) *InstallationStatus {
	lower := strings.ToLower(line)

	status := &InstallationStatus{
		Formula: pkg,
	}

	switch {
	case strings.Contains(lower, "downloading"):
		status.Stage = "downloading"
		status.Progress = 25
	case strings.Contains(lower, "installing"):
		status.Stage = "installing"
		status.Progress = 50
	case strings.Contains(lower, "pouring"):
		status.Stage = "installing"
		status.Progress = 60
	case strings.Contains(lower, "linking"):
		status.Stage = "linking"
		status.Progress = 90
	case strings.Contains(lower, "installed"):
		status.Stage = "completed"
		status.Progress = 100
	default:
		return nil
	}

	return status
}

// ExecuteCommand executes a raw brew command (fallback)
func (c *Client) ExecuteCommand(ctx context.Context, args []string) error {
	//nolint:gosec // brewPath is validated at client creation, args are user commands
	cmd := exec.CommandContext(ctx, c.brewPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
