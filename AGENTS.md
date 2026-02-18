# AGENTS.md (zenswitch)

## このディレクトリ配下の方針
- 実装言語は `Go` を使用する。
- CLIアプリとして `go run ./cmd/zen` で実行できる形を維持する。
- 仕様変更時は、まずロジック層をテストで担保してからCLI層を更新する。

## 現在のCLI仕様
- `zen` は、許可リストに含まれない起動中アプリを終了する（macOS専用）。
- 許可リストの既定値は `Terminal` / `iTerm2` / `Ghostty` / `Finder` / `Dock` / `System Settings` / `Activity Monitor`。
- `zen`（CLI自身）は常に終了対象外。
- `zen list` で、オプション適用後の最終許可リストを表示して終了する（アプリ終了は行わない）。
- `zen add "<App>"` で、許可リストへ永続追加する。
- `zen remove "<App>"` で、許可リストから永続除外する。
- `-h` / `--help` / `zen help` でヘルプを表示する。
- `--dry-run` で、終了対象のみ表示して実際には終了しない。
- `--allow "App1,App2"` で既定許可リストに追加できる。
- `--allow-only --allow "App1,App2"` で既定許可リストを使わず指定したアプリのみを許可できる。
- `--disallow "App1,App2"` で許可リスト（既定 + `--allow`）から除外できる。
- `--allow-only` は `--allow` とセット必須。
- 設定ファイルは既定で `~/.config/zenswitch/config.json` を参照し、`--config` で任意パスを指定できる。
- 設定ファイルキーは `replaceDefaultAllowed` / `allowedApps` / `disallowedApps`。

## 実行コマンド
- フォーマット: `gofmt -w ./cmd ./internal`
- テスト: `go test ./...`
- 実行: `go run ./cmd/zen`
- 許可リスト表示: `go run ./cmd/zen list`
- 許可リストへ追加: `go run ./cmd/zen add Visual Studio Code`
- 許可リストから除外: `go run ./cmd/zen remove Ghostty`
- ヘルプ表示: `go run ./cmd/zen -h`
- 終了対象の確認のみ: `go run ./cmd/zen --dry-run`
- 許可リスト追加: `go run ./cmd/zen --allow "Ghostty,Visual Studio Code"`
- 許可リスト置換: `go run ./cmd/zen --allow-only --allow "Ghostty,Visual Studio Code"`
- 許可リスト除外: `go run ./cmd/zen --disallow "Ghostty,Visual Studio Code"`
- 設定ファイル指定: `go run ./cmd/zen --config "/path/to/config.json" --list`
