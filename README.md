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

ユーザー指定の許可リストを使う場合:

```bash
# 現在の有効な許可リストを表示（アプリは終了しない）
zen --list

# 既定の許可リストに追加
zen --allow "Ghostty,Visual Studio Code"

# 指定したアプリのみを許可（既定リストは使わない）
zen --allow-only --allow "Ghostty,Visual Studio Code"

# 許可リストから除外（既定リストにも追加リストにも適用）
zen --disallow "Ghostty,Visual Studio Code"

# オプション適用後の最終許可リストを表示
zen --list --allow "Arc" --disallow "Ghostty"
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
