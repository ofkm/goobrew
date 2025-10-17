// Package version provides version information for goobrew.
// Version values are typically set via ldflags during the build process.
package version

// Build information variables set via ldflags during compilation.
var (
	// Version is the semantic version number (e.g., "0.1.2").
	Version = "0.2.0"
	// Commit is the git commit hash of the build.
	Commit = "none"
	// BuildTime is the timestamp when the binary was built.
	BuildTime = "unknown"
)

// GetVersion returns a formatted version string including the program name.
func GetVersion() string {
	return "goobrew version " + Version
}

// GetFullVersion returns a detailed version string including the version number,
// git commit hash, and build timestamp. Each component is displayed on a separate line.
func GetFullVersion() string {
	return GetVersion() + "\n" +
		"commit: " + Commit + "\n" +
		"built: " + BuildTime
}
