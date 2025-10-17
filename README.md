# goobrew 🍺 (like goober)

[![Tests](https://github.com/ofkm/goobrew/actions/workflows/test.yml/badge.svg)](https://github.com/ofkm/goobrew/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ofkm/goobrew)](https://goreportcard.com/report/github.com/ofkm/goobrew)
[![License](https://img.shields.io/github/license/ofkm/goobrew)](LICENSE)

> [!NOTE]
> This is a WIP eventually want to make this a standlone brew client instead of a wrapper.

A opinionated fast and friendly wrapper for Homebrew written in Go.

## Features

- **Fast**: Built in Go for optimal performance
- **User-friendly**: Enhanced UX with nerd-font icon indicators and timing information
- **Full Homebrew compatibility**: All standard brew commands work seamlessly
- **Convenient aliases**: Shorter commands for common operations

## Installation

### Using go install

```bash
go install github.com/ofkm/goobrew@latest
```

### From source

```bash
git clone https://github.com/ofkm/goobrew.git
cd goobrew
make build
# Or: make install to install to $GOPATH/bin
```

## Usage

goobrew wraps Homebrew commands with enhanced feedback and user experience:

### Common Commands

```bash
# Install packages (alias: i)
goobrew install wget
goobrew i git

# Uninstall packages (aliases: remove, rm)
goobrew uninstall wget
goobrew rm git

# Search for packages (alias: s)
goobrew search python
goobrew s node

# Update Homebrew (alias: up)
goobrew update

# Upgrade packages
goobrew upgrade
goobrew upgrade git

# List installed packages (alias: ls)
goobrew list
goobrew ls

# Show package information
goobrew info git

# Show version
goobrew version
```

### Pass-through to Homebrew

Any command not explicitly handled by goobrew is automatically passed through to Homebrew:

```bash
goobrew doctor
goobrew cleanup
goobrew services list
```

## Why goobrew?

- **Better feedback**: Clear emoji indicators and execution timing
- **Familiar syntax**: Works exactly like brew, but better
- **Shorter commands**: Convenient aliases for frequently used commands
- **No learning curve**: If you know brew, you know goobrew

## Requirements

- Go 1.21 or higher
- Homebrew installed on your system

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
