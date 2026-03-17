//go:build ignore

package main

import (
	"fmt"
	"os"
	"strings"
)

// Generate SVG terminal screenshots for README
func main() {
	generateChartView()
	generateTableView()
	fmt.Println("Generated docs/chart-view.svg and docs/table-view.svg")
}

func generateChartView() {
	lines := []string{
		"\033[1;38;5;205m★ gr-stars\033[0m",
		"",
		"\033[1;38;5;212;48;5;236m Claude / Anthropic \033[0m  \033[38;5;240mAI / LLM\033[0m  \033[38;5;240mGo\033[0m  \033[38;5;240mWeb\033[0m",
		"",
	}

	type repo struct {
		name  string
		stars string
		bar   int
		color string
	}

	repos := []repo{
		{"anthropics/claude-code", "25.3k", 50, "#ff0000"},
		{"modelcontextprotocol/s..", "16.8k", 33, "#ff5f00"},
		{"anthropics/courses", "12.1k", 24, "#ff8700"},
		{"anthropics/anthropic-..", "10.5k", 21, "#ffaf00"},
		{"langgenius/dify", "9.2k", 18, "#ffd700"},
		{"punkpeye/awesome-mcp-..", "8.7k", 17, "#ffff00"},
		{"mckaywrigley/chatbot-ui", "7.9k", 15, "#afff00"},
		{"vercel/ai-chatbot", "7.1k", 14, "#5fff00"},
		{"supabase/supabase", "6.4k", 12, "#00ff00"},
		{"cline/cline", "5.8k", 11, "#00ff5f"},
	}

	for i, r := range repos {
		cursor := "  "
		if i == 0 {
			cursor = "▸ "
		}
		name := fmt.Sprintf("%-28s", r.name)
		stars := fmt.Sprintf("%8s", r.stars)
		bar := strings.Repeat("█", r.bar)
		lines = append(lines, fmt.Sprintf("%s%s %s %s", cursor, name, stars, bar))
	}

	lines = append(lines, "")
	lines = append(lines, "\033[38;5;241mj/k: select  Enter: open  Tab: category  v: view (chart)  /: search  q: quit\033[0m")

	writeSVG("docs/chart-view.svg", lines, "gr-stars — Chart View")
}

func generateTableView() {
	lines := []string{
		"\033[1;38;5;205m★ gr-stars\033[0m",
		"",
		"\033[38;5;240mClaude / Anthropic\033[0m  \033[1;38;5;212;48;5;236m AI / LLM \033[0m  \033[38;5;240mGo\033[0m  \033[38;5;240mWeb\033[0m",
		"",
		"\033[1m #   Repository                       Stars  Language     Description\033[0m",
		"─────────────────────────────────────────────────────────────────────────────────",
	}

	type row struct {
		rank int
		name string
		stars string
		lang string
		desc string
	}

	rows := []row{
		{1, "tensorflow/tensorflow", "187.2k", "C++", "An Open Source Machine Learning..."},
		{2, "pytorch/pytorch", "86.4k", "Python", "Tensors and Dynamic neural net..."},
		{3, "huggingface/transformers", "142.1k", "Python", "Transformers: State-of-the-art..."},
		{4, "openai/openai-python", "25.8k", "Python", "The official Python library for..."},
		{5, "langchain-ai/langchain", "102.3k", "Python", "Build context-aware reasoning..."},
		{6, "ggerganov/llama.cpp", "75.6k", "C++", "LLM inference in C/C++"},
		{7, "ollama/ollama", "120.4k", "Go", "Get up and running with Llama..."},
		{8, "deepseek-ai/DeepSeek-V3", "30.1k", "Python", "DeepSeek-V3 model implementat..."},
	}

	for i, r := range rows {
		cursor := "  "
		if i == 2 {
			cursor = "▸ "
		}
		lines = append(lines, fmt.Sprintf("%s%2d  %-30s %8s  %-12s %s", cursor, r.rank, r.name, r.stars, r.lang, r.desc))
	}

	lines = append(lines, "")
	lines = append(lines, "\033[38;5;241mj/k: select  Enter: open  Tab: category  v: view (table)  /: search  q: quit\033[0m")

	writeSVG("docs/table-view.svg", lines, "gr-stars — Table View")
}

func writeSVG(filename string, lines []string, title string) {
	lineHeight := 20
	paddingTop := 50
	paddingBottom := 20
	paddingX := 20
	width := 820
	height := paddingTop + len(lines)*lineHeight + paddingBottom

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">
  <defs>
    <style>
      @import url('https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;700&amp;display=swap');
      .term { font-family: 'JetBrains Mono', 'SF Mono', 'Menlo', monospace; font-size: 13px; }
      .term-bold { font-family: 'JetBrains Mono', 'SF Mono', 'Menlo', monospace; font-size: 13px; font-weight: bold; }
    </style>
  </defs>
  <!-- Window chrome -->
  <rect width="%d" height="%d" rx="8" fill="#1a1b26"/>
  <rect x="0" y="0" width="%d" height="36" rx="8" fill="#24283b"/>
  <rect x="0" y="28" width="%d" height="8" fill="#24283b"/>
  <circle cx="18" cy="18" r="6" fill="#ff5f57"/>
  <circle cx="38" cy="18" r="6" fill="#febc2e"/>
  <circle cx="58" cy="18" r="6" fill="#28c840"/>
  <text x="%d" y="22" text-anchor="middle" fill="#a9b1d6" class="term" font-size="12">%s</text>
`, width, height, width, height, width, height, width, width, width/2, title))

	for i, line := range lines {
		y := paddingTop + i*lineHeight
		rendered := ansiToSVG(line, paddingX, y)
		sb.WriteString(rendered)
	}

	sb.WriteString("</svg>\n")

	if err := os.WriteFile(filename, []byte(sb.String()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", filename, err)
		os.Exit(1)
	}
}

type textSpan struct {
	text  string
	color string
	bold  bool
	bg    string
}

func ansiToSVG(line string, x, y int) string {
	spans := parseANSI(line)
	var sb strings.Builder
	cx := x
	charWidth := 7.8

	for _, span := range spans {
		if span.text == "" {
			continue
		}

		textLen := float64(len(span.text)) * charWidth

		if span.bg != "" {
			sb.WriteString(fmt.Sprintf(`  <rect x="%.1f" y="%d" width="%.1f" height="18" fill="%s"/>
`, float64(cx), y-13, textLen, span.bg))
		}

		class := "term"
		if span.bold {
			class = "term-bold"
		}
		color := span.color
		if color == "" {
			color = "#a9b1d6"
		}

		escaped := strings.ReplaceAll(span.text, "<", "&lt;")
		escaped = strings.ReplaceAll(escaped, ">", "&gt;")
		escaped = strings.ReplaceAll(escaped, "&", "&amp;")

		sb.WriteString(fmt.Sprintf(`  <text x="%.1f" y="%d" fill="%s" class="%s">%s</text>
`, float64(cx), y, color, class, escaped))

		cx += int(textLen)
	}
	return sb.String()
}

func parseANSI(s string) []textSpan {
	var spans []textSpan
	currentColor := ""
	currentBold := false
	currentBg := ""
	i := 0
	textStart := 0

	ansi256ToHex := map[string]string{
		"205": "#ff5faf", "212": "#ff87d7", "240": "#585858",
		"236": "#303030", "117": "#87d7ff", "228": "#ffff87",
		"246": "#949494", "242": "#6c6c6c", "241": "#626262",
		"196": "#ff0000", "202": "#ff5f00", "208": "#ff8700",
		"214": "#ffaf00", "220": "#ffd700", "226": "#ffff00",
		"190": "#d7ff00", "154": "#afff00", "118": "#87ff00",
		"82":  "#5fff00", "46":  "#00ff00", "47":  "#00ff5f",
		"48":  "#00ff87", "49":  "#00ffaf", "50":  "#00ffd7",
		"51":  "#00ffff", "45":  "#00d7ff", "39":  "#00afff",
		"33":  "#0087ff", "27":  "#005fff", "255": "#eeeeee",
	}

	for i < len(s) {
		if s[i] == '\033' && i+1 < len(s) && s[i+1] == '[' {
			if textStart < i {
				spans = append(spans, textSpan{
					text: s[textStart:i], color: currentColor,
					bold: currentBold, bg: currentBg,
				})
			}

			end := i + 2
			for end < len(s) && s[end] != 'm' {
				end++
			}
			if end < len(s) {
				codes := strings.Split(s[i+2:end], ";")
				for j := 0; j < len(codes); j++ {
					switch codes[j] {
					case "0":
						currentColor = ""
						currentBold = false
						currentBg = ""
					case "1":
						currentBold = true
					case "38":
						if j+2 < len(codes) && codes[j+1] == "5" {
							if hex, ok := ansi256ToHex[codes[j+2]]; ok {
								currentColor = hex
							}
							j += 2
						}
					case "48":
						if j+2 < len(codes) && codes[j+1] == "5" {
							if hex, ok := ansi256ToHex[codes[j+2]]; ok {
								currentBg = hex
							}
							j += 2
						}
					}
				}
				i = end + 1
				textStart = i
			} else {
				i++
			}
		} else {
			i++
		}
	}

	if textStart < len(s) {
		spans = append(spans, textSpan{
			text: s[textStart:], color: currentColor,
			bold: currentBold, bg: currentBg,
		})
	}

	return spans
}
