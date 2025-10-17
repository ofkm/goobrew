package version

// Version information set via ldflags during build
var (
	// Version is the semantic version
	Version = "0.1.1"
	// Commit is the git commit hash
	Commit = "none"
	// BuildTime is when the binary was built
	BuildTime = "unknown"
)

// GetVersion returns the formatted version string
func GetVersion() string {
	return "goobrew version " + Version
}

// GetFullVersion returns version with all details
func GetFullVersion() string {
	return GetVersion() + "\n" +
		"commit: " + Commit + "\n" +
		"built: " + BuildTime
}
