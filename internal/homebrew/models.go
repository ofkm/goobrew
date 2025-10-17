package homebrew

import (
	"encoding/json"
	"time"
)

// Formula represents a Homebrew formula
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

type Versions struct {
	Stable string `json:"stable"`
	Head   string `json:"head,omitempty"`
	Bottle bool   `json:"bottle"`
}

type URLs struct {
	Stable StableURL `json:"stable,omitempty"`
	Head   HeadURL   `json:"head,omitempty"`
}

type StableURL struct {
	URL      string `json:"url"`
	Tag      string `json:"tag,omitempty"`
	Revision string `json:"revision,omitempty"`
	Using    string `json:"using,omitempty"`
	Checksum string `json:"checksum,omitempty"`
}

type HeadURL struct {
	URL    string `json:"url"`
	Branch string `json:"branch,omitempty"`
	Using  string `json:"using,omitempty"`
}

type Bottle struct {
	Rebuild int                   `json:"rebuild"`
	RootURL string                `json:"root_url"`
	Files   map[string]BottleFile `json:"files"`
}

type BottleFile struct {
	Cellar string `json:"cellar"`
	URL    string `json:"url"`
	Sha256 string `json:"sha256"`
}

type KegOnlyReason struct {
	Reason      string `json:"reason"`
	Explanation string `json:"explanation"`
}

type Option struct {
	Option      string `json:"option"`
	Description string `json:"description"`
}

type Requirement struct {
	Name     string   `json:"name"`
	Cask     string   `json:"cask,omitempty"`
	Download string   `json:"download,omitempty"`
	Version  string   `json:"version,omitempty"`
	Contexts []string `json:"contexts"`
	Specs    []string `json:"specs"`
}

type InstalledInfo struct {
	Version               string       `json:"version"`
	UsedOptions           []string     `json:"used_options"`
	BuiltAsBottle         bool         `json:"built_as_bottle"`
	PouredFromBottle      bool         `json:"poured_from_bottle"`
	Time                  int64        `json:"time"`
	RuntimeDependencies   []Dependency `json:"runtime_dependencies"`
	InstalledAsDependency bool         `json:"installed_as_dependency"`
	InstalledOnRequest    bool         `json:"installed_on_request"`
}

type Dependency struct {
	FullName         string `json:"full_name"`
	Version          string `json:"version"`
	Revision         int    `json:"revision"`
	PkgVersion       string `json:"pkg_version"`
	DeclaredDirectly bool   `json:"declared_directly"`
}

type Service struct {
	Name       string          `json:"name"`
	RunType    string          `json:"run_type"`
	RunAtLoad  bool            `json:"run_at_load,omitempty"`
	KeepAlive  json.RawMessage `json:"keep_alive,omitempty"` // Can be bool or object
	WorkingDir string          `json:"working_dir,omitempty"`
}

// GetKeepAliveBool tries to extract a boolean value from KeepAlive
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

type RubyChecksum struct {
	Sha256 string `json:"sha256"`
}

// SearchResult represents search results
type SearchResult struct {
	Formulae []string `json:"formulae"`
	Casks    []string `json:"casks"`
}

// InstallationStatus represents installation progress
type InstallationStatus struct {
	Formula   string
	Stage     string // "downloading", "installing", "linking", "completed", "failed"
	Progress  int    // 0-100
	StartTime time.Time
	Error     error
}
