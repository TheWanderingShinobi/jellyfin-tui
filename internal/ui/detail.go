package ui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/TheWanderingShinobi/jellyfin-tui/internal/errors"
	"github.com/TheWanderingShinobi/jellyfin-tui/internal/jellyfin"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	detailTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4")).
				Padding(0, 1)

	detailInfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			MarginTop(1)

	detailActionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				MarginTop(1)
)

type detailModel struct {
	item   *jellyfin.MediaItem
	client *jellyfin.Client
}

func newDetailModel(client *jellyfin.Client) detailModel {
	return detailModel{
		client: client,
	}
}

func (m detailModel) Init() tea.Cmd {
	return nil
}

func (m detailModel) Update(msg tea.Msg) (detailModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, m.playMedia
		case "p":
			return m, m.addToPlaylist
		case "esc", "q":
			return m, m.back
		}
	case jellyfin.MediaItem:
		m.item = &msg
	}
	return m, nil
}

func (m detailModel) View() string {
	if m.item == nil {
		return "Loading..."
	}

	var b strings.Builder

	b.WriteString(detailTitleStyle.Render(m.item.Name))
	b.WriteString("\n\n")

	b.WriteString(detailInfoStyle.Render(fmt.Sprintf("Type: %s\n", m.item.Type)))
	if m.item.Overview != "" {
		b.WriteString(detailInfoStyle.Render(fmt.Sprintf("Overview: %s\n", m.item.Overview)))
	}
	if m.item.CommunityRating > 0 {
		b.WriteString(detailInfoStyle.Render(fmt.Sprintf("Rating: %.1f\n", m.item.CommunityRating)))
	}

	b.WriteString("\n")
	b.WriteString(detailActionStyle.Render("Press Enter to play"))
	b.WriteString("\n")
	b.WriteString(detailActionStyle.Render("Press 'p' to add to playlist"))
	b.WriteString("\n")
	b.WriteString(detailActionStyle.Render("Press 'q' or Esc to go back"))

	return b.String()
}

func (m detailModel) playMedia() tea.Msg {
	if m.item != nil {
		streamURL := m.client.GetStreamURL(m.item.ID)
		err := playWithMPV(streamURL)
		if err != nil {
			return errors.NewAPIError(fmt.Sprintf("Failed to play media: %v", err))
		}
	}
	return nil
}

func (m detailModel) addToPlaylist() tea.Msg {
	if m.item != nil {
		return addToPlaylistMsg{itemID: m.item.ID}
	}
	return nil
}

func (m detailModel) back() tea.Msg {
	return showBrowseMsg{}
}

func playWithMPV(streamURL string) error {
	cmd := exec.Command("mpv", streamURL)
	return cmd.Run()
}

type addToPlaylistMsg struct {
	itemID string
}
