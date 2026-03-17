package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/engineer-fumi/gr-stars/tui"
)

func main() {
	query := flag.String("query", "", "Custom search query")
	category := flag.String("category", "", "Category name to start with")
	flag.Parse()

	_ = category // reserved for future use

	m := tui.NewModel(*query)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
