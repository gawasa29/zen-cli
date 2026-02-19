# zen-cli 🧘 - Close distracting apps in one command

zen-cli is a macOS-only focus CLI. Run `zen` to quit running apps that are not in your allow-list, so you can reset your desktop before deep work. It is designed for local-only operation with predictable defaults, and requires Go 1.22+ only when you build from source.

## Install

### Requirements
- macOS (`osascript` and `pkill` available)
- Homebrew (recommended install path)
- Go 1.22+ (only for source builds)

### Package manager
```bash
brew tap gawasa29/tap
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

---

# zen-cli 🧘 - 1コマンドで集中を邪魔するアプリを閉じる

zen-cli は macOS 専用の集中用 CLI です。`zen` を実行すると、許可リストに含まれない起動中アプリを終了し、深い作業に入る前にデスクトップを素早く整えられます。動作はローカル完結で、ソースビルド時のみ Go 1.22+ が必要です。

## Install

### Requirements
- macOS（`osascript` と `pkill` が利用可能）
- Homebrew（推奨インストール経路）
- Go 1.22+（ソースからビルドする場合のみ）

### Package manager
```bash
brew tap gawasa29/tap
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

1. 使えるコマンドを確認します。
2. 何が終了対象になるかを事前確認します。
3. 実際に終了処理を実行します。
4. 永続許可リストを更新します。

```bash
zen --help
zen --dry-run
zen
zen list
zen add "Visual Studio Code"
zen remove "Ghostty"
```

## Features

- 許可対象外のフォアグラウンドアプリを 1 コマンドで終了します。
- `zen add` と `zen remove` で永続許可リストを管理できます。
- `--allow`、`--allow-only`、`--disallow` で一時的な実行条件を上書きできます。
- `zen list`、`--list`、`--dry-run` で安全に事前確認できます。
- CLI 自身（`zen`）は常に終了対象から除外されます。

## Commands

- `zen`: 有効な許可リストに含まれないアプリを終了します。
- `zen list`: 有効な許可リストを表示して終了します。
- `zen add APP_NAME`: アプリ名を許可リスト設定に追加します。
- `zen remove APP_NAME`: アプリ名を許可リスト設定から除外します。
- `zen help [list|add|remove]`: ルートコマンドまたはサブコマンドのヘルプを表示します。

## Configuration

既定の設定ファイルパス:
- `XDG_CONFIG_HOME` が設定されている場合は `$XDG_CONFIG_HOME/zen-cli/config.json`。
- それ以外は `~/.config/zen-cli/config.json`。

設定ファイル形式:

```json
{
  "replaceDefaultAllowed": false,
  "allowedApps": ["Ghostty", "Visual Studio Code"],
  "disallowedApps": ["Slack"]
}
```

任意の設定ファイルを使う例:

```bash
zen --config "/path/to/config.json" --list
zen --config "/path/to/config.json" --dry-run
```

## Docs

- [プロジェクト方針と運用](AGENTS.md)
- [CLI エントリポイント](cmd/zen/main.go)
- [コアのアプリ判定ロジック](internal/zencli/zencli.go)
- [CLI テスト](cmd/zen/main_test.go)
- [コアロジックのテスト](internal/zencli/zencli_test.go)

## Privacy / Permissions / Limitations

- Privacy: 外部サービスへの送信は行わず、処理はすべてローカルで完結します。
- Permissions: `osascript` で対象アプリを制御するため、macOS のオートメーション権限が必要になる場合があります。
- Limitations: macOS 専用であり、保存前のデータがある状態で終了すると内容が失われる可能性があります。

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

- [Projects monorepo](https://github.com/gawasa29/Projects) - ソースリポジトリ。
- [Issue tracker](https://github.com/gawasa29/Projects/issues) - バグ報告と機能要望。

## License

MIT（`LICENSE`）
