# Release Process

This document explains how to create and publish releases for goobrew.

## Versioning

goobrew follows [Semantic Versioning](https://semver.org/):

- MAJOR version for incompatible API changes
- MINOR version for new functionality in a backwards compatible manner
- PATCH version for backwards compatible bug fixes

## Creating a Release

### 1. Update Version in Code

**IMPORTANT**: Before creating a tag, update the version number in the code:

1. Edit `internal/version/version.go` and update the `Version` constant:
   ```go
   Version = "X.Y.Z"  // Change from previous version
   ```

2. Commit this change:
   ```bash
   git add internal/version/version.go
   git commit -m "Bump version to vX.Y.Z"
   git push origin main
   ```

### 2. Ensure Quality

1. Ensure all tests pass:

   ```bash
   make test
   make lint
   ```

2. Verify CI is green on GitHub Actions

### 3. Create and Push the Tag

1. Create a git tag following the format `vX.Y.Z`:

   ```bash
   git tag -a v0.2.0 -m "Release v0.2.0: Add new features"
   ```

2. Push the tag to GitHub:
   ```bash
   git push origin v0.2.0
   ```

3. The GitHub Actions workflow will automatically:
   - Build binaries for multiple platforms (Linux, macOS, Windows)
   - Create checksums for all binaries
   - Generate release notes
   - Create a GitHub Release with all artifacts
   - Test the installation

### 4. Verify the Release

Once the tag is pushed and the workflow completes:

1. Check the GitHub Actions workflow succeeded
2. Verify the release appears at https://github.com/ofkm/goobrew/releases
3. Test installation:

```bash
go install github.com/ofkm/goobrew@v0.2.0
```

Or the latest version:

```bash
go install github.com/ofkm/goobrew@latest
```

## Version Information

The version is set in `internal/version/version.go` and **must be updated before each release**.

For builds with additional metadata (commit, build time), the Makefile and CI use ldflags to inject:
- Git commit hash
- Build timestamp

Users can check the version with:

```bash
goobrew version
```

## Automated Release Workflow

The `.github/workflows/release.yml` automatically runs when you push a tag and:

1. Builds binaries for Linux, macOS, Windows (amd64 + arm64)
2. Embeds full version info (version, commit, build time)
3. Creates SHA256 checksums
4. Generates release notes from commits
5. Publishes GitHub Release with all artifacts
6. Tests the installation

## Local Development

For local development builds:

```bash
make build          # Build with version info
make install        # Install to $GOPATH/bin
```

The Makefile automatically:

- Extracts version from git tags
- Includes the git commit hash
- Adds the build timestamp

## First Release Checklist

For the first release (v0.1.0):

- [ ] All tests passing
- [ ] README.md is complete and accurate
- [ ] License file is present
- [ ] go.mod has correct module path
- [ ] Create git tag: `git tag -a v0.1.0 -m "Initial release"`
- [ ] Push tag: `git push origin v0.1.0`
- [ ] Verify installation: `go install github.com/ofkm/goobrew@v0.1.0`
- [ ] Create GitHub release (optional but recommended)

## GitHub Releases (Optional)

While not required for `go install` to work, creating GitHub releases provides a nice UI:

1. Go to https://github.com/ofkm/goobrew/releases/new
2. Select the tag you just created
3. Add release notes describing the changes
4. Publish the release

## Example Release Commands

```bash
# Tag and release v0.1.0
git tag -a v0.1.0 -m "Initial release with basic functionality"
git push origin v0.1.0

# Tag and release v0.2.0 with new features
git tag -a v0.2.0 -m "Add search caching and parallel operations"
git push origin v0.2.0

# Tag and release v0.2.1 with bug fixes
git tag -a v0.2.1 -m "Fix race condition in cache loading"
git push origin v0.2.1
```

## Troubleshooting

### `go install` says "no matching versions"

Make sure:

1. The tag is pushed to GitHub: `git push origin vX.Y.Z`
2. The tag follows semantic versioning format: `vX.Y.Z`
3. Wait a minute for Go's proxy cache to update: https://proxy.golang.org/github.com/ofkm/goobrew/@v/list

### Version shows as "dev"

If you build without using `make` or `go install`, the version will show as "dev". Use:

```bash
make build    # For local builds with version info
make install  # To install with version info
```
