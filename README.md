# ZenSwitch

ZenSwitch は macOS 専用の CLI アプリです。`zen` を実行すると、許可リストに含まれない起動中アプリを終了します。

## 概要
- 集中前に不要アプリを一括で閉じる。
- 許可リストは `zen add` / `zen remove` で永続管理できる。
- `zen list` / `zen --dry-run` で実行前確認ができる。

## 対応環境
- macOS のみ（`osascript` / `pkill` が利用できる環境）
- Go 1.22 以上（ソースからビルドする場合）

## インストール
```bash
cd zenswitch
go build -o zen ./cmd/zen
sudo install -m 755 zen /usr/local/bin/zen
```

## クイックスタート
```bash
# 許可リスト外のアプリを終了
zen

# 現在の許可リストを確認
zen list

# 許可リストに追加（永続化）
zen add Visual Studio Code

# 許可リストから除外（永続化）
zen remove Ghostty

# 終了対象の確認のみ（実際には終了しない）
zen --dry-run

# ヘルプ
zen -h
```

## コマンド
- `zen`: 許可リスト外のアプリを終了
- `zen list`: 有効な許可リストを表示
- `zen add APP_NAME [APP_NAME ...]`: 許可リストへ追加して設定ファイルに保存
- `zen remove APP_NAME [APP_NAME ...]`: 許可リストから除外して設定ファイルに保存
- `zen help [list|add|remove]`: ヘルプ表示

## 主要オプション
- `--dry-run`: 終了対象だけ表示して終了（アプリ終了しない）
- `--config PATH`: 使用する設定ファイルを指定
- `--allow APP1,APP2`: 今回実行に限って許可リストへ追加
- `--allow-only`: 既定許可リストを使わず、明示的な許可のみを使用
- `--disallow APP1,APP2`: 今回実行に限って許可リストから除外
- `--list`: 許可リスト表示（互換のため維持）
- `-h`, `--help`: ヘルプ表示

## 設定ファイル
既定パス:
- `XDG_CONFIG_HOME` がある場合: `$XDG_CONFIG_HOME/zenswitch/config.json`
- ない場合: `~/.config/zenswitch/config.json`

形式:
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

## 既定の許可リスト
- Terminal
- iTerm2
- Ghostty
- Finder
- Dock
- System Settings
- Activity Monitor
- zen（CLI 自身）

## 終了コード
- `0`: 正常終了
- `1`: 実行エラー（引数不正、設定読み込み失敗、アプリ終了失敗など）
- `2`: 非対応OS（macOS以外）

## 安全上の注意
- 保存していないデータがあるアプリは、実行前に保存してください。
- 初回利用時は `zen --dry-run` で対象確認してから `zen` 実行を推奨します。
- 初回実行時は macOS のオートメーション権限許可が必要になる場合があります。

## テスト
```bash
go test ./...
```

## リリース手順
```bash
mkdir -p dist
GOOS=darwin GOARCH=arm64 go build -o dist/zen-macos-arm64 ./cmd/zen
GOOS=darwin GOARCH=amd64 go build -o dist/zen-macos-amd64 ./cmd/zen
shasum -a 256 dist/zen-macos-arm64 dist/zen-macos-amd64 > dist/checksums.txt
```

## サポート
- バグ報告・要望は GitHub Issues: `https://github.com/gawasa29/Projects/issues`
- 報告時は `macOSバージョン`、`実行コマンド`、`エラーメッセージ` を含めてください。

## ライセンス
- MIT License (`LICENSE`)
