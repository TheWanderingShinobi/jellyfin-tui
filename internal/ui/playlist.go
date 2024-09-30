package ui

import (
	"fmt"
	"strings"

	"github.com/TheWanderingShinobi/jellyfin-tui/internal/jellyfin"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	playlistTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4")).
				Padding(0, 1)

	playlistItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA"))

	playlistSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Background(lipgloss.Color("#FAFAFA"))
)

type playlistModel struct {
	playlists    []jellyfin.Playlist
	cursor       int
	selectedItem string
	client       *jellyfin.Client
}

func newPlaylistModel(client *jellyfin.Client) playlistModel {
	return playlistModel{
		client: client,
	}
}

func (m playlistModel) Init() tea.Cmd {
	return m.fetchPlaylists
}

func (m playlistModel) Update(msg tea.Msg) (playlistModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.playlists)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.playlists) > 0 {
				return m, m.addToPlaylist
			}
		case "n":
			return m, m.createNewPlaylist
		case "esc", "q":
			return m, m.back
		}
	case playlistsMsg:
		m.playlists = msg.playlists
	}
	return m, nil
}

func (m playlistModel) View() string {
	var b strings.Builder

	b.WriteString(playlistTitleStyle.Render("Playlists"))
	b.WriteString("\n\n")

	for i, playlist := range m.playlists {
		if i == m.cursor {
			b.WriteString(playlistSelectedStyle.Render("> " + playlist.Name))
		} else {
			b.WriteString(playlistItemStyle.Render("  " + playlist.Name))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString("Press Enter to add item to selected playlist\n")
	b.WriteString("Press 'n' to create a new playlist\n")
	b.WriteString("Press 'q' or Esc to go back")

	return b.String()
}

func (m playlistModel) fetchPlaylists() tea.Msg {
	playlists, err := m.client.GetPlaylists()
	if err != nil {
		return errorMsg{err}
	}
	return playlistsMsg{playlists: playlists}
}

func (m playlistModel) addToPlaylist() tea.Msg {
	if m.cursor < len(m.playlists) {
		err := m.client.AddToPlaylist(m.playlists[m.cursor].ID, m.selectedItem)
		if err != nil {
			return errorMsg{err}
		}
		return playlistUpdateMsg{message: "Item added to playlist"}
	}
	return nil
}

func (m playlistModel) createNewPlaylist() tea.Cmd {
	return func() tea.Msg {
		return showCreatePlaylistMsg{}
	}
}

func (m playlistModel) back() tea.Msg {
	return showBrowseMsg{}
}

type playlistsMsg struct {
	playlists []jellyfin.Playlist
}

type playlistUpdateMsg struct {
	message string
}

type showCreatePlaylistMsg struct{}
