package tui

import (
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbletea"
	"github.com/engineer-fumi/gr-stars/github"
)

type ViewMode int

const (
	ViewChart ViewMode = iota
	ViewTable
)

type Model struct {
	categories   []Category
	catIndex     int
	repos        []github.Repository
	loading      bool
	err          error
	width        int
	height       int
	viewMode     ViewMode
	cursor       int
	customQuery  string
	inputMode    bool
	inputBuffer  string
}

type reposMsg struct {
	repos []github.Repository
	err   error
}

func NewModel(initialQuery string) Model {
	cats := Categories
	if initialQuery != "" {
		cats = append(cats, Category{Name: "Custom", Query: initialQuery})
	}
	return Model{
		categories: cats,
		loading:    true,
		viewMode:   ViewChart,
	}
}

func (m Model) Init() tea.Cmd {
	return m.fetchRepos()
}

func (m Model) fetchRepos() tea.Cmd {
	query := m.categories[m.catIndex].Query
	return func() tea.Msg {
		repos, err := github.SearchRepositories(query)
		return reposMsg{repos: repos, err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.inputMode {
			return m.handleInput(msg)
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.catIndex = (m.catIndex + 1) % len(m.categories)
			m.loading = true
			return m, m.fetchRepos()
		case "shift+tab":
			m.catIndex = (m.catIndex - 1 + len(m.categories)) % len(m.categories)
			m.loading = true
			return m, m.fetchRepos()
		case "j", "down":
			if len(m.repos) > 0 {
				m.cursor = (m.cursor + 1) % len(m.repos)
			}
			return m, nil
		case "k", "up":
			if len(m.repos) > 0 {
				m.cursor = (m.cursor - 1 + len(m.repos)) % len(m.repos)
			}
			return m, nil
		case "enter":
			if len(m.repos) > 0 && m.cursor < len(m.repos) {
				return m, openBrowser(m.repos[m.cursor].HTMLURL)
			}
			return m, nil
		case "v":
			if m.viewMode == ViewChart {
				m.viewMode = ViewTable
			} else {
				m.viewMode = ViewChart
			}
			return m, nil
		case "/":
			m.inputMode = true
			m.inputBuffer = ""
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case reposMsg:
		m.loading = false
		m.repos = msg.repos
		m.err = msg.err
		m.cursor = 0
		return m, nil
	}
	return m, nil
}

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

func (m Model) handleInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.inputMode = false
		if m.inputBuffer != "" {
			found := false
			for i, c := range m.categories {
				if c.Name == "Custom" {
					m.categories[i].Query = m.inputBuffer
					m.catIndex = i
					found = true
					break
				}
			}
			if !found {
				m.categories = append(m.categories, Category{Name: "Custom", Query: m.inputBuffer})
				m.catIndex = len(m.categories) - 1
			}
			m.loading = true
			return m, m.fetchRepos()
		}
		return m, nil
	case "esc":
		m.inputMode = false
		m.inputBuffer = ""
		return m, nil
	case "backspace":
		if len(m.inputBuffer) > 0 {
			m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
		}
		return m, nil
	default:
		if len(msg.String()) == 1 || msg.String() == " " {
			m.inputBuffer += msg.String()
		}
		return m, nil
	}
}

func (m Model) View() string {
	return renderView(m)
}
