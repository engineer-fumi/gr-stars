# ★ gr-stars

一款在终端中可视化 GitHub 热门仓库的 CLI 工具。
支持分类筛选、柱状图/表格视图切换，以及直接在浏览器中打开仓库。

**[English](README.md)** | **[日本語](README.ja.md)**

<p align="center">
  <img src="docs/chart-view.svg" alt="Chart View" width="780">
</p>

<p align="center">
  <img src="docs/table-view.svg" alt="Table View" width="780">
</p>

## 功能特性

- **柱状图视图** — 通过彩色柱状图直观比较 Star 数量
- **表格视图** — 以列表形式展示仓库名称、Star 数、编程语言和描述
- **分类筛选** — 在 Claude/Anthropic、AI/LLM、Go、Web 等分类间切换
- **自定义搜索** — 按 `/` 键输入任意 GitHub 搜索查询
- **浏览器打开** — 按 Enter 在浏览器中打开选中的仓库
- **跨平台** — 支持 macOS、Linux 和 Windows

## 安装

```bash
go install github.com/engineer-fumi/gr-stars@latest
```

或从源码构建：

```bash
git clone https://github.com/engineer-fumi/gr-stars.git
cd gr-stars
go build -o gr-stars .
```

## 使用方法

```bash
# 启动（默认分类：Claude / Anthropic）
gr-stars

# 使用自定义搜索查询启动
gr-stars -query "topic:rust stars:>5000"
```

### GitHub Token（推荐）

设置 Personal Access Token 以避免 GitHub API 速率限制：

```bash
export GITHUB_TOKEN="ghp_your_token_here"
```

> 不设置 Token 也可以使用，但频繁请求可能会触发速率限制。

## 快捷键

| 按键 | 操作 |
|------|------|
| `j` / `↓` | 选择下一个仓库 |
| `k` / `↑` | 选择上一个仓库 |
| `Enter` | 在浏览器中打开选中的仓库 |
| `Tab` | 切换到下一个分类 |
| `Shift+Tab` | 切换到上一个分类 |
| `v` | 切换柱状图 ⇄ 表格视图 |
| `/` | 进入自定义搜索模式 |
| `Esc` | 退出搜索模式 |
| `q` | 退出 |

## 分类

| 分类 | 查询 |
|------|------|
| Claude / Anthropic | `claude OR anthropic OR mcp-server topic:claude` |
| AI / LLM | `llm OR ai OR gpt topic:machine-learning` |
| Go | `language:go stars:>1000` |
| Web | `topic:react OR topic:vue OR topic:nextjs` |

## 技术栈

- [Go](https://go.dev/)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI 框架
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — 终端样式
- [GitHub Search API](https://docs.github.com/en/rest/search/search) — 仓库搜索

## 许可证

[MIT](LICENSE)
