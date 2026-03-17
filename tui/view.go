package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	catStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	catActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			Background(lipgloss.Color("236")).
			Padding(0, 1)

	repoNameStyle = lipgloss.NewStyle().
			Width(30).
			Foreground(lipgloss.Color("117"))

	starsStyle = lipgloss.NewStyle().
			Width(8).
			Align(lipgloss.Right).
			Foreground(lipgloss.Color("228"))

	langStyle = lipgloss.NewStyle().
			Width(12).
			Foreground(lipgloss.Color("246"))

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("242"))

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("255")).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))
)

func renderView(m Model) string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("★ gr-stars"))
	b.WriteString("\n")

	// Category tabs
	var tabs []string
	for i, c := range m.categories {
		if i == m.catIndex {
			tabs = append(tabs, catActiveStyle.Render(c.Name))
		} else {
			tabs = append(tabs, catStyle.Render(c.Name))
		}
	}
	b.WriteString(strings.Join(tabs, "  "))
	b.WriteString("\n\n")

	// Input mode
	if m.inputMode {
		b.WriteString(inputStyle.Render("Search: " + m.inputBuffer + "█"))
		b.WriteString("\n\n")
	}

	// Loading / Error
	if m.loading {
		b.WriteString("Loading...\n")
		b.WriteString(renderHelp(m))
		return b.String()
	}
	if m.err != nil {
		b.WriteString(fmt.Sprintf("Error: %v\n", m.err))
		b.WriteString(renderHelp(m))
		return b.String()
	}
	if len(m.repos) == 0 {
		b.WriteString("No repositories found.\n")
		b.WriteString(renderHelp(m))
		return b.String()
	}

	// Render based on view mode
	switch m.viewMode {
	case ViewChart:
		b.WriteString(renderChart(m))
	case ViewTable:
		b.WriteString(renderTable(m))
	}

	b.WriteString(renderHelp(m))
	return b.String()
}

func renderChart(m Model) string {
	var b strings.Builder

	maxStars := m.repos[0].Stars
	barMaxWidth := m.width - 45
	if barMaxWidth < 20 {
		barMaxWidth = 20
	}

	for i, repo := range m.repos {
		name := repo.FullName
		if len(name) > 28 {
			name = name[:28]
		}

		barWidth := 0
		if maxStars > 0 {
			barWidth = (repo.Stars * barMaxWidth) / maxStars
		}
		if barWidth < 1 {
			barWidth = 1
		}

		bar := strings.Repeat("█", barWidth)

		// Color gradient based on rank
		color := rankColor(i)
		barStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))

		cursor := "  "
		if i == m.cursor {
			cursor = "▸ "
		}

		line := fmt.Sprintf("%s%s %s %s",
			cursorStyle(i == m.cursor).Render(cursor),
			repoNameStyle.Render(name),
			starsStyle.Render(formatStars(repo.Stars)),
			barStyle.Render(bar),
		)
		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}

func renderTable(m Model) string {
	var b strings.Builder

	// Header
	header := fmt.Sprintf("%s %s %s %s %s",
		lipgloss.NewStyle().Width(4).Bold(true).Render("#"),
		lipgloss.NewStyle().Width(30).Bold(true).Render("Repository"),
		lipgloss.NewStyle().Width(10).Align(lipgloss.Right).Bold(true).Render("Stars"),
		lipgloss.NewStyle().Width(12).Bold(true).Render("Language"),
		lipgloss.NewStyle().Bold(true).Render("Description"),
	)
	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n")

	for i, repo := range m.repos {
		name := repo.FullName
		if len(name) > 28 {
			name = name[:28]
		}

		desc := repo.Description
		descMax := m.width - 64
		if descMax < 10 {
			descMax = 10
		}
		if len(desc) > descMax {
			desc = desc[:descMax-3] + "..."
		}

		lang := repo.Language
		if lang == "" {
			lang = "-"
		}

		cursor := "  "
		if i == m.cursor {
			cursor = "▸ "
		}

		line := fmt.Sprintf("%s%s %s %s %s %s",
			cursorStyle(i == m.cursor).Render(cursor),
			lipgloss.NewStyle().Width(4).Render(fmt.Sprintf("%2d", i+1)),
			repoNameStyle.Render(name),
			starsStyle.Render(formatStars(repo.Stars)),
			langStyle.Render(lang),
			descStyle.Render(desc),
		)
		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}

func cursorStyle(selected bool) lipgloss.Style {
	if selected {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	}
	return lipgloss.NewStyle()
}

func renderHelp(m Model) string {
	mode := "chart"
	if m.viewMode == ViewTable {
		mode = "table"
	}
	return helpStyle.Render(fmt.Sprintf(
		"j/k: select  Enter: open  Tab: category  v: view (%s)  /: search  q: quit",
		mode,
	))
}

func formatStars(n int) string {
	if n >= 1000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}

func rankColor(i int) string {
	colors := []string{
		"196", "202", "208", "214", "220",
		"226", "190", "154", "118", "82",
		"46", "47", "48", "49", "50",
		"51", "45", "39", "33", "27",
	}
	if i < len(colors) {
		return colors[i]
	}
	return "244"
}
