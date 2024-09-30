package ui

import (
	"fmt"
	"strings"

	"github.com/TheWanderingShinobi/jellyfin-tui/internal/jellyfin"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	searchPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4")).
				Padding(0, 1)

	searchResultStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA"))

	searchSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Background(lipgloss.Color("#FAFAFA"))
)

type searchModel struct {
	query   string
	results []jellyfin.MediaItem
	cursor  int
	client  *jellyfin.Client
}

func newSearchModel(client *jellyfin.Client) searchModel {
	return searchModel{
		client: client,
	}
}

func (m searchModel) Init() tea.Cmd {
	return nil
}

func (m searchModel) Update(msg tea.Msg) (searchModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.query != "" {
				return m, m.search
			} else if len(m.results) > 0 {
				return m, m.selectItem
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.results)-1 {
				m.cursor++
			}
		case "esc", "q":
			return m, m.back
		case "backspace":
			if len(m.query) > 0 {
				m.query = m.query[:len(m.query)-1]
			}
		default:
			m.query += msg.String()
		}
	case searchResultMsg:
		m.results = msg.results
		m.cursor = 0
	}
	return m, nil
}

func (m searchModel) View() string {
	var b strings.Builder

	b.WriteString(searchPromptStyle.Render("Search: " + m.query))
	b.WriteString("\n\n")

	if len(m.results) > 0 {
		for i, item := range m.results {
			result := fmt.Sprintf("%s (%s)", item.Name, item.Type)
			if i == m.cursor {
				b.WriteString(searchSelectedStyle.Render(result))
			} else {
				b.WriteString(searchResultStyle.Render(result))
			}
			b.WriteString("\n")
		}
	} else if m.query != "" {
		b.WriteString("No results found.")
	}

	b.WriteString("\n")
	b.WriteString("Press Enter to search or select item\n")
	b.WriteString("Press 'q' or Esc to go back")

	return b.String()
}

func (m searchModel) search() tea.Msg {
	results, err := m.client.Search(m.query)
	if err != nil {
		return errorMsg{err}
	}
	return searchResultMsg{results: results}
}

func (m searchModel) selectItem() tea.Msg {
	if len(m.results) > 0 {
		return showDetailMsg{item: m.results[m.cursor]}
	}
	return nil
}

func (m searchModel) back() tea.Msg {
	return showBrowseMsg{}
}

type searchResultMsg struct {
	results []jellyfin.MediaItem
}
