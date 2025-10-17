# goobrew 🍺

A fast and UX-friendly wrapper for Homebrew written in Go.

## Features

- 🚀 **Fast**: Built in Go for optimal performance
- 😊 **User-friendly**: Enhanced UX with emoji indicators and timing information
- 🔄 **Full Homebrew compatibility**: All standard brew commands work seamlessly
- ⚡ **Convenient aliases**: Shorter commands for common operations

## Installation

```bash
go install github.com/ofkm/goobrew@latest
```

Or build from source:

```bash
git clone https://github.com/ofkm/goobrew.git
cd goobrew
go build -o goobrew
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
