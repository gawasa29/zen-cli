# ZenSwitch

ZenSwitch は macOS 専用の CLI アプリです。`zen` を実行すると、あらかじめ許可したアプリ（Terminal など）を除く起動中アプリを終了します。

## 目的
- 集中作業に入る前に、不要なアプリを一括で閉じる。

## セットアップ
```bash
cd zenswitch
go build -o zen ./cmd/zen
```

必要なら `PATH` の通った場所へ配置してください。

```bash
mv zen /usr/local/bin/zen
```

## 実行方法
```bash
zen
```

シンプルな操作コマンド:

```bash
# 許可リスト外のアプリを終了
zen

# 現在の有効な許可リストを表示（終了はしない）
zen list

# 許可リストへ追加（永続化）
zen add Visual Studio Code

# 許可リストから除外（永続化）
zen remove Ghostty

# ヘルプ
zen -h
zen help
zen add -h
```

設定ファイル（既定: `~/.config/zenswitch/config.json`）でも許可リストを指定できます。

```json
{
  "replaceDefaultAllowed": false,
  "allowedApps": ["Ghostty", "Visual Studio Code"],
  "disallowedApps": ["Slack"]
}
```

任意の設定ファイルを使う場合:

```bash
zen --config "/path/to/config.json" --list
zen --config "/path/to/config.json" --dry-run
```

高度なオプション（必要な場合のみ）:

```bash
# 終了対象の確認のみ（実際には終了しない）
zen --dry-run

# 既定の許可リストに追加
zen --allow "Ghostty,Visual Studio Code"

# 指定したアプリのみを許可（既定リストは使わない）
zen --allow-only --allow "Ghostty,Visual Studio Code"

# 許可リストから除外（既定リストにも追加リストにも適用）
zen --disallow "Ghostty,Visual Studio Code"
```

終了対象から除外する既定アプリは以下です。
- Terminal
- iTerm2
- Ghostty
- Finder
- Dock
- System Settings
- Activity Monitor
- zen（CLI 自身）

## テスト方法
```bash
go test ./...
```

## リリース方法
```bash
# バイナリ作成
GOOS=darwin GOARCH=arm64 go build -o dist/zen-macos-arm64 ./cmd/zen
GOOS=darwin GOARCH=amd64 go build -o dist/zen-macos-amd64 ./cmd/zen
```

必要に応じてチェックサムを生成し、配布してください。
