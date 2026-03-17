# ★ gr-stars

GitHub で人気のリポジトリをターミナル上で可視化する CLI ツールです。
カテゴリ別フィルタ、バーチャート / テーブルの表示切替、ブラウザでの直接オープンに対応しています。

<p align="center">
  <img src="docs/chart-view.svg" alt="Chart View" width="780">
</p>

<p align="center">
  <img src="docs/table-view.svg" alt="Table View" width="780">
</p>

## Features

- **バーチャート表示** — スター数を色付きバーで直感的に比較
- **テーブル表示** — リポジトリ名・スター数・言語・説明を一覧表示
- **カテゴリフィルタ** — Claude / Anthropic、AI / LLM、Go、Web など関心ごとに切替
- **カスタム検索** — `/` キーで自由な GitHub 検索クエリを入力
- **ブラウザオープン** — 選択したリポジトリを Enter で直接ブラウザに開く
- **クロスプラットフォーム** — macOS / Linux / Windows 対応

## Installation

```bash
go install github.com/engineer-fumi/gr-stars@latest
```

または、ソースからビルド：

```bash
git clone https://github.com/engineer-fumi/gr-stars.git
cd gr-stars
go build -o gr-stars .
```

## Usage

```bash
# 基本起動（デフォルトカテゴリ: Claude / Anthropic）
gr-stars

# カスタムクエリを指定して起動
gr-stars -query "topic:rust stars:>5000"
```

### GitHub Token（推奨）

GitHub API のレート制限を緩和するために、Personal Access Token を設定できます：

```bash
export GITHUB_TOKEN="ghp_your_token_here"
```

> トークンなしでも動作しますが、短時間に多くのリクエストを送ると制限に達する場合があります。

## Keybindings

| Key | Action |
|-----|--------|
| `j` / `↓` | 次のリポジトリを選択 |
| `k` / `↑` | 前のリポジトリを選択 |
| `Enter` | 選択中のリポジトリをブラウザで開く |
| `Tab` | 次のカテゴリへ切替 |
| `Shift+Tab` | 前のカテゴリへ切替 |
| `v` | チャート ⇄ テーブル表示を切替 |
| `/` | カスタム検索モードに入る |
| `Esc` | 検索モードを終了 |
| `q` | 終了 |

## Categories

| Category | Query |
|----------|-------|
| Claude / Anthropic | `claude OR anthropic OR mcp-server topic:claude` |
| AI / LLM | `llm OR ai OR gpt topic:machine-learning` |
| Go | `language:go stars:>1000` |
| Web | `topic:react OR topic:vue OR topic:nextjs` |

## Tech Stack

- [Go](https://go.dev/)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI フレームワーク
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — ターミナルスタイリング
- [GitHub Search API](https://docs.github.com/en/rest/search/search) — リポジトリ検索

## License

MIT
