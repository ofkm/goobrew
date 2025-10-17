package homebrew

import (
	"encoding/json"
	"time"
)

// Formula represents a Homebrew formula with all its metadata.
// It contains information about package name, version, dependencies,
// installation status, and various other attributes from Homebrew's JSON API.
type Formula struct {
	Name                 string              `json:"name"`
	FullName             string              `json:"full_name"`
	Tap                  string              `json:"tap"`
	OldName              string              `json:"oldname,omitempty"`
	Aliases              []string            `json:"aliases"`
	VersionedFormulae    []string            `json:"versioned_formulae"`
	Desc                 string              `json:"desc"`
	License              string              `json:"license,omitempty"`
	Homepage             string              `json:"homepage"`
	Versions             Versions            `json:"versions"`
	Urls                 URLs                `json:"urls,omitempty"`
	Revision             int                 `json:"revision"`
	VersionScheme        int                 `json:"version_scheme"`
	Bottle               Bottle              `json:"bottle"`
	KegOnly              bool                `json:"keg_only"`
	KegOnlyReason        *KegOnlyReason      `json:"keg_only_reason,omitempty"`
	Options              []Option            `json:"options"`
	BuildDependencies    []string            `json:"build_dependencies"`
	Dependencies         []string            `json:"dependencies"`
	TestDependencies     []string            `json:"test_dependencies"`
	RecommendedDeps      []string            `json:"recommended_dependencies"`
	OptionalDeps         []string            `json:"optional_dependencies"`
	UsesFromMacos        []json.RawMessage   `json:"uses_from_macos"`
	UsesFromMacosBounds  []map[string]string `json:"uses_from_macos_bounds"`
	Requirements         []Requirement       `json:"requirements"`
	ConflictsWith        []string            `json:"conflicts_with"`
	ConflictsWithReasons []string            `json:"conflicts_with_reasons"`
	LinkOverwrite        []string            `json:"link_overwrite"`
	Caveats              string              `json:"caveats,omitempty"`
	Installed            []InstalledInfo     `json:"installed"`
	LinkedKeg            string              `json:"linked_keg,omitempty"`
	Pinned               bool                `json:"pinned"`
	Outdated             bool                `json:"outdated"`
	Deprecated           bool                `json:"deprecated"`
	DeprecationDate      string              `json:"deprecation_date,omitempty"`
	DeprecationReason    string              `json:"deprecation_reason,omitempty"`
	Disabled             bool                `json:"disabled"`
	DisableDate          string              `json:"disable_date,omitempty"`
	DisableReason        string              `json:"disable_reason,omitempty"`
	PostInstallDefined   bool                `json:"post_install_defined"`
	Service              *Service            `json:"service,omitempty"`
	TapGitHead           string              `json:"tap_git_head,omitempty"`
	RubySourcePath       string              `json:"ruby_source_path,omitempty"`
	RubySourceChecksum   RubyChecksum        `json:"ruby_source_checksum,omitempty"`
}

// Versions contains version information for a formula.
type Versions struct {
	Stable string `json:"stable"`          // Stable is the stable version number
	Head   string `json:"head,omitempty"`  // Head is the HEAD version if available
	Bottle bool   `json:"bottle"`          // Bottle indicates if precompiled bottles are available
}

// URLs contains download URLs for a formula.
type URLs struct {
	Stable StableURL `json:"stable,omitempty"` // Stable is the stable release URL
	Head   HeadURL   `json:"head,omitempty"`   // Head is the development HEAD URL
}

// StableURL contains the stable release download information.
type StableURL struct {
	URL      string `json:"url"`                // URL is the download URL
	Tag      string `json:"tag,omitempty"`      // Tag is the git tag if applicable
	Revision string `json:"revision,omitempty"` // Revision is the git revision
	Using    string `json:"using,omitempty"`    // Using specifies download strategy
	Checksum string `json:"checksum,omitempty"` // Checksum is the file checksum
}

// HeadURL contains the development HEAD download information.
type HeadURL struct {
	URL    string `json:"url"`               // URL is the repository URL
	Branch string `json:"branch,omitempty"`  // Branch is the git branch
	Using  string `json:"using,omitempty"`   // Using specifies download strategy
}

// Bottle contains precompiled binary package information.
type Bottle struct {
	Rebuild int                   `json:"rebuild"` // Rebuild number
	RootURL string                `json:"root_url"` // RootURL is the base URL for bottles
	Files   map[string]BottleFile `json:"files"`    // Files maps platform to bottle file info
}

// BottleFile contains information about a specific bottle file.
type BottleFile struct {
	Cellar string `json:"cellar"` // Cellar path
	URL    string `json:"url"`    // URL is the download URL
	Sha256 string `json:"sha256"` // Sha256 is the file checksum
}

// KegOnlyReason explains why a formula is keg-only.
type KegOnlyReason struct {
	Reason      string `json:"reason"`      // Reason is a short reason code
	Explanation string `json:"explanation"` // Explanation is a detailed explanation
}

// Option represents an installation option for a formula.
type Option struct {
	Option      string `json:"option"`      // Option is the option flag
	Description string `json:"description"` // Description explains the option
}

// Requirement represents a system requirement for a formula.
type Requirement struct {
	Name     string   `json:"name"`               // Name is the requirement name
	Cask     string   `json:"cask,omitempty"`     // Cask specifies a required cask
	Download string   `json:"download,omitempty"` // Download is a download URL
	Version  string   `json:"version,omitempty"`  // Version is the required version
	Contexts []string `json:"contexts"`           // Contexts where this applies
	Specs    []string `json:"specs"`              // Specs where this applies
}

// InstalledInfo contains information about an installed version of a formula.
type InstalledInfo struct {
	Version               string       `json:"version"`
	UsedOptions           []string     `json:"used_options"`
	BuiltAsBottle         bool         `json:"built_as_bottle"`
	PouredFromBottle      bool         `json:"poured_from_bottle"`
	Time                  int64        `json:"time"`
	RuntimeDependencies   []Dependency `json:"runtime_dependencies"`
	InstalledAsDependency bool         `json:"installed_as_dependency"`
	InstalledOnRequest    bool         `json:"installed_on_request"`    // InstalledOnRequest indicates if installed explicitly
}

// Dependency represents a formula dependency.
type Dependency struct {
	FullName         string `json:"full_name"`         // FullName is the full formula name
	Version          string `json:"version"`           // Version is the dependency version
	Revision         int    `json:"revision"`          // Revision is the formula revision
	PkgVersion       string `json:"pkg_version"`       // PkgVersion is version including revision
	DeclaredDirectly bool   `json:"declared_directly"` // DeclaredDirectly indicates if explicitly declared
}

// Service defines a background service configuration for a formula.
type Service struct {
	Name       string          `json:"name"`                    // Name is the service name
	RunType    string          `json:"run_type"`                // RunType defines when the service runs
	RunAtLoad  bool            `json:"run_at_load,omitempty"`   // RunAtLoad indicates if service starts at load
	KeepAlive  json.RawMessage `json:"keep_alive,omitempty"`    // KeepAlive can be bool or object
	WorkingDir string          `json:"working_dir,omitempty"`   // WorkingDir is the service working directory
}

// GetKeepAliveBool attempts to extract a boolean value from the KeepAlive field.
// The KeepAlive field can be either a boolean or an object. This method returns
// true if KeepAlive is explicitly true or if it's an object (considered "enabled").
// Returns false if KeepAlive is empty or explicitly false.
func (s *Service) GetKeepAliveBool() bool {
	if len(s.KeepAlive) == 0 {
		return false
	}

	// Try to unmarshal as bool
	var boolVal bool
	if err := json.Unmarshal(s.KeepAlive, &boolVal); err == nil {
		return boolVal
	}

	// If it's an object, consider it as "enabled"
	return true
}

// RubyChecksum contains checksum information for the formula's Ruby source.
type RubyChecksum struct {
	Sha256 string `json:"sha256"` // Sha256 is the SHA-256 checksum
}

// SearchResult contains the results of a package search operation.
type SearchResult struct {
	Formulae []string `json:"formulae"` // Formulae lists matching formula names
	Casks    []string `json:"casks"`    // Casks lists matching cask names
}

// InstallationStatus represents the current state of a package installation.
// It is used to communicate progress updates from the Install method to callers.
type InstallationStatus struct {
	Formula   string        // Formula is the package name being installed
	Stage     string        // Stage is one of: "downloading", "installing", "linking", "completed", "failed"
	Progress  int           // Progress is a percentage from 0-100
	StartTime time.Time     // StartTime is when the installation began
	Error     error         // Error contains any error that occurred
}
