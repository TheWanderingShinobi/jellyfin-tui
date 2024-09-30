package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	helpTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	helpSectionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true)

	helpContentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA"))
)

type helpModel struct{}

func newHelpModel() helpModel {
	return helpModel{}
}

func (m helpModel) Init() tea.Cmd {
	return nil
}

func (m helpModel) Update(msg tea.Msg) (helpModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, m.back
		}
	}
	return m, nil
}

func (m helpModel) View() string {
	var b strings.Builder

	b.WriteString(helpTitleStyle.Render("Jellyfin TUI Help"))
	b.WriteString("\n\n")

	b.WriteString(helpSectionStyle.Render("Navigation"))
	b.WriteString("\n")
	b.WriteString(helpContentStyle.Render("j/down: Move cursor down\n"))
	b.WriteString(helpContentStyle.Render("k/up: Move cursor up\n"))
	b.WriteString(helpContentStyle.Render("enter: Select item\n"))
	b.WriteString(helpContentStyle.Render("q/esc: Go back/quit\n"))
	b.WriteString("\n")

	b.WriteString(helpSectionStyle.Render("Browse View"))
	b.WriteString("\n")
	b.WriteString(helpContentStyle.Render("f: Toggle filters\n"))
	b.WriteString(helpContentStyle.Render("s: Search\n"))
	b.WriteString(helpContentStyle.Render("n: Next page\n"))
	b.WriteString(helpContentStyle.Render("p: Previous page\n"))
	b.WriteString("\n")

	b.WriteString(helpSectionStyle.Render("Detail View"))
	b.WriteString("\n")
	b.WriteString(helpContentStyle.Render("enter: Play media\n"))
	b.WriteString(helpContentStyle.Render("p: Add to playlist\n"))
	b.WriteString("\n")

	b.WriteString(helpSectionStyle.Render("Playlist View"))
	b.WriteString("\n")
	b.WriteString(helpContentStyle.Render("enter: Add item to selected playlist\n"))
	b.WriteString(helpContentStyle.Render("n: Create new playlist\n"))
	b.WriteString("\n")

	b.WriteString(helpSectionStyle.Render("Settings View"))
	b.WriteString("\n")
	b.WriteString(helpContentStyle.Render("enter: Edit selected setting\n"))
	b.WriteString("\n")

	b.WriteString(helpContentStyle.Render("Press 'q' or Esc to go back"))

	return b.String()
}

func (m helpModel) back() tea.Msg {
	return showBrowseMsg{}
}
