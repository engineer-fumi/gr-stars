---
title: "Go + Bubble Tea で GitHub の人気リポジトリを可視化する CLI ツールを作った"
emoji: "⭐"
type: "tech"
topics: ["go", "cli", "github", "bubbletea", "tui"]
published: false
---

## はじめに

「今、GitHub で何が流行っているんだろう？」

ブラウザで GitHub Trending を開けばいい話ですが、ターミナルで作業しているときにわざわざブラウザに切り替えるのは面倒です。そこで、**ターミナル上で GitHub の人気リポジトリをサクッと確認できる CLI ツール**を作りました。

https://github.com/engineer-fumi/gr-stars

![Demo](https://raw.githubusercontent.com/engineer-fumi/gr-stars/main/docs/demo.gif)

## できること

- **バーチャート表示**: スター数を色付きバーで直感的に比較
- **テーブル表示**: リポジトリ名・スター数・言語・説明を一覧表示
- **カテゴリフィルタ**: Claude / Anthropic、AI / LLM、Go、Web など関心ごとに切替
- **カスタム検索**: `/` キーで自由に GitHub 検索クエリを入力
- **ブラウザオープン**: 気になるリポジトリを `Enter` で即座にブラウザで開く

## 技術スタック

| ライブラリ | 用途 |
|-----------|------|
| [Bubble Tea](https://github.com/charmbracelet/bubbletea) | TUI フレームワーク（Elm Architecture） |
| [Lip Gloss](https://github.com/charmbracelet/lipgloss) | ターミナルのスタイリング |
| GitHub Search API | リポジトリ検索 |

## アーキテクチャ

```
gr-stars/
├── main.go              # エントリポイント
├── github/
│   └── client.go        # GitHub API クライアント
└── tui/
    ├── categories.go    # 検索カテゴリ定義
    ├── model.go         # Elm Architecture の Model / Update
    └── view.go          # View（描画処理）
```

Bubble Tea は **Elm Architecture** を採用しています。シンプルに `Model`（状態）→ `Update`（状態変更）→ `View`（描画）のサイクルで動きます。

## 実装のポイント

### 1. GitHub Search API を叩く

GitHub の Search API はトークンなしでも使えますが、レート制限があります。`GITHUB_TOKEN` 環境変数があれば自動でヘッダーに付与するようにしました。

```go
func SearchRepositories(query string) ([]Repository, error) {
    u := fmt.Sprintf(
        "https://api.github.com/search/repositories?q=%s&sort=stars&order=desc&per_page=20",
        url.QueryEscape(query),
    )

    req, _ := http.NewRequest("GET", u, nil)
    req.Header.Set("Accept", "application/vnd.github.v3+json")
    if token := os.Getenv("GITHUB_TOKEN"); token != "" {
        req.Header.Set("Authorization", "Bearer "+token)
    }
    // ...
}
```

### 2. Bubble Tea の非同期コマンド

API 呼び出しは `tea.Cmd` として非同期実行します。これにより UI がブロックされず、ローディング表示ができます。

```go
func (m Model) fetchRepos() tea.Cmd {
    query := m.categories[m.catIndex].Query
    return func() tea.Msg {
        repos, err := github.SearchRepositories(query)
        return reposMsg{repos: repos, err: err}
    }
}
```

### 3. バーチャートの色グラデーション

ランキング順に赤→黄→緑→青のグラデーションをつけています。ANSI 256 色を活用することで、ターミナルでもリッチな表現が可能です。

```go
func rankColor(i int) string {
    colors := []string{
        "196", "202", "208", "214", "220", // 赤〜黄
        "226", "190", "154", "118", "82",  // 黄〜緑
        "46", "47", "48", "49", "50",      // 緑〜シアン
        "51", "45", "39", "33", "27",      // シアン〜青
    }
    if i < len(colors) {
        return colors[i]
    }
    return "244"
}
```

### 4. クロスプラットフォームなブラウザオープン

`runtime.GOOS` で OS を判定し、各プラットフォームのコマンドを使い分けます。

```go
func openBrowser(url string) tea.Cmd {
    return func() tea.Msg {
        var cmd *exec.Cmd
        switch runtime.GOOS {
        case "darwin":
            cmd = exec.Command("open", url)
        case "linux":
            cmd = exec.Command("xdg-open", url)
        default:
            cmd = exec.Command("cmd", "/c", "start", url)
        }
        _ = cmd.Start()
        return nil
    }
}
```

### 5. カテゴリ定義の拡張性

カテゴリは単純な構造体のスライスで定義しているため、追加が容易です。

```go
var Categories = []Category{
    {Name: "Claude / Anthropic", Query: "claude OR anthropic OR mcp-server topic:claude"},
    {Name: "AI / LLM",          Query: "llm OR ai OR gpt topic:machine-learning"},
    {Name: "Go",                Query: "language:go stars:>1000"},
    {Name: "Web",               Query: "topic:react OR topic:vue OR topic:nextjs"},
}
```

GitHub の検索構文がそのまま使えるので、`language:rust stars:>5000` のような高度なクエリもカスタム検索で対応できます。

## 使い方

### インストール

```bash
go install github.com/engineer-fumi/gr-stars@latest
```

### 起動

```bash
# GitHub Token を設定（推奨）
export GITHUB_TOKEN="ghp_your_token_here"

# 起動
gr-stars
```

### キーバインド

| キー | 操作 |
|------|------|
| `j` / `↓` | 次のリポジトリを選択 |
| `k` / `↑` | 前のリポジトリを選択 |
| `Enter` | ブラウザで開く |
| `Tab` | 次のカテゴリ |
| `v` | チャート ⇄ テーブル切替 |
| `/` | カスタム検索 |
| `q` | 終了 |

## Bubble Tea で TUI を作ってみた感想

Charm 社のエコシステム（Bubble Tea + Lip Gloss）は非常に完成度が高く、Go で TUI を作るなら現状最良の選択肢だと感じました。

**良かった点：**
- Elm Architecture のおかげで状態管理が明確
- `tea.Cmd` による非同期処理が直感的
- Lip Gloss のスタイリングが CSS ライクで書きやすい

**注意点：**
- ターミナルの幅に応じたレイアウト調整は自前で実装する必要がある
- 日本語などマルチバイト文字の幅計算には注意が必要

## まとめ

Go + Bubble Tea を使えば、100行程度のコードでリッチな TUI アプリケーションが構築できます。ターミナルでの作業が多い方は、ぜひ試してみてください。

PR や Issue も歓迎です！

https://github.com/engineer-fumi/gr-stars
