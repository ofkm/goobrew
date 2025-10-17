package version

// Version information set via ldflags during build
var (
	// Version is the semantic version
	Version = "dev"
	// Commit is the git commit hash
	Commit = "none"
	// BuildTime is when the binary was built
	BuildTime = "unknown"
)

// GetVersion returns the formatted version string
func GetVersion() string {
	if Version == "dev" {
		return "goobrew development version"
	}
	return "goobrew version " + Version
}

// GetFullVersion returns version with all details
func GetFullVersion() string {
	return GetVersion() + "\n" +
		"commit: " + Commit + "\n" +
		"built: " + BuildTime
}
