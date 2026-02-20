[English](README.md) | [日本語](README.ja.md)

<p align="center">
  <img src="assets/zen-cli-hero-20260220.png" alt="zen-cli logo" width="480">
</p>

# zen-cli 🧘 - Close distracting apps in one command

zen-cli is a macOS-only focus CLI. Run `zen` to quit running apps that are not in your allow-list, so you can reset your desktop before deep work. It is designed for local-only operation with predictable defaults, and requires Go 1.22+ only when you build from source.

## Install

### Requirements
- macOS (`osascript` and `pkill` available)
- Homebrew (recommended install path)
- Go 1.22+ (only for source builds)

### Package manager
```bash
brew install gawasa29/tap/zen-cli
```

### Build from source
```bash
git clone https://github.com/gawasa29/Projects.git
cd Projects/zen-cli
go build -o zen ./cmd/zen
sudo install -m 755 zen /usr/local/bin/zen
```

## Quick Start

1. Confirm available commands.
2. Preview what would be closed.
3. Run the cleanup.
4. Manage the persistent allow-list.

```bash
zen --help
zen --dry-run
zen
zen list
zen add "Visual Studio Code"
zen remove "Ghostty"
```

## Features

- Quits non-allowed foreground apps in one command.
- Keeps a persistent allow-list with `zen add` and `zen remove`.
- Supports one-shot overrides with `--allow`, `--allow-only`, and `--disallow`.
- Provides safe preview modes with `zen list`, `--list`, and `--dry-run`.
- Always excludes the CLI process itself (`zen`) from quit targets.

## Commands

- `zen`: Quit apps outside the effective allow-list.
- `zen list`: Print the effective allow-list and exit.
- `zen add APP_NAME`: Add app names to allow-list config.
- `zen remove APP_NAME`: Remove app names from allow-list config.
- `zen help [list|add|remove]`: Show help for root command or a subcommand.

## Configuration

Default config path:
- `$XDG_CONFIG_HOME/zen-cli/config.json` when `XDG_CONFIG_HOME` is set.
- `~/.config/zen-cli/config.json` otherwise.

Config shape:

```json
{
  "replaceDefaultAllowed": false,
  "allowedApps": ["Ghostty", "Visual Studio Code"],
  "disallowedApps": ["Slack"]
}
```

Use a custom config path:

```bash
zen --config "/path/to/config.json" --list
zen --config "/path/to/config.json" --dry-run
```

## Docs

- [Project policy and workflow](AGENTS.md)
- [CLI entry point](cmd/zen/main.go)
- [Core app filtering logic](internal/zencli/zencli.go)
- [CLI tests](cmd/zen/main_test.go)
- [Core logic tests](internal/zencli/zencli_test.go)

## Privacy / Permissions / Limitations

- Privacy: zen-cli does not send data to external services; all processing is local.
- Permissions: macOS may request Automation permission so `osascript` can control target apps.
- Limitations: macOS-only, and quitting apps can discard unsaved work if you do not save first.

## Getting started (dev)

```bash
cd zen-cli
go test ./...
go run ./cmd/zen --dry-run
go run ./cmd/zen list
```

## Build from source

```bash
cd zen-cli
go build -o zen ./cmd/zen
./zen --help
```

## Release

```bash
cd zen-cli
mkdir -p dist
GOOS=darwin GOARCH=arm64 go build -o dist/zen-macos-arm64 ./cmd/zen
GOOS=darwin GOARCH=amd64 go build -o dist/zen-macos-amd64 ./cmd/zen
shasum -a 256 dist/zen-macos-arm64 dist/zen-macos-amd64 > dist/checksums.txt
```

## Related

- [Projects monorepo](https://github.com/gawasa29/Projects) - source repository.
- [Issue tracker](https://github.com/gawasa29/Projects/issues) - bug reports and feature requests.

## License

MIT (`LICENSE`)
