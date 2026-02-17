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

終了対象から除外する既定アプリは以下です。
- Terminal
- iTerm2
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
