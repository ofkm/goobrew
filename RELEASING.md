# Release Process

This document explains how to create and publish releases for goobrew.

## Versioning

goobrew follows [Semantic Versioning](https://semver.org/):
- MAJOR version for incompatible API changes
- MINOR version for new functionality in a backwards compatible manner  
- PATCH version for backwards compatible bug fixes

## Creating a Release

### 1. Prepare the Release

1. Ensure all tests pass:
   ```bash
   make test
   make lint
   ```

2. Update CHANGELOG.md (if you have one) with the changes in this release

3. Commit any pending changes:
   ```bash
   git add .
   git commit -m "Prepare for vX.Y.Z release"
   ```

### 2. Create and Push the Tag

1. Create a git tag following the format `vX.Y.Z`:
   ```bash
   git tag -a v0.1.0 -m "Release v0.1.0"
   ```

2. Push the tag to GitHub:
   ```bash
   git push origin v0.1.0
   ```

### 3. Verify the Release

Once the tag is pushed, users can install it with:

```bash
go install github.com/ofkm/goobrew@v0.1.0
```

Or the latest version:

```bash
go install github.com/ofkm/goobrew@latest
```

## Version Information

The version information is embedded in the binary at build time using `ldflags`. Users can check the version with:

```bash
goobrew version
```

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
