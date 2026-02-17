# AGENTS.md (zenswitch)

## このディレクトリ配下の方針
- 実装言語は `Go` を使用する。
- CLIアプリとして `go run ./cmd/zen` で実行できる形を維持する。
- 仕様変更時は、まずロジック層をテストで担保してからCLI層を更新する。

## 実行コマンド
- フォーマット: `gofmt -w ./cmd ./internal`
- テスト: `go test ./...`
- 実行: `go run ./cmd/zen`
