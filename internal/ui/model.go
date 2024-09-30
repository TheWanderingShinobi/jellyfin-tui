package ui

import (
	"github.com/TheWanderingShinobi/jellyfin-tui/internal/config"
	"github.com/TheWanderingShinobi/jellyfin-tui/internal/errors"
	"github.com/TheWanderingShinobi/jellyfin-tui/internal/jellyfin"
	"github.com/charmbracelet/bubbletea"
)

type Model struct {
	client        *jellyfin.Client
	config        config.Config
	state         string
	loginModel    loginModel
	browseModel   browseModel
	detailModel   detailModel
	searchModel   searchModel
	playlistModel playlistModel
	settingsModel settingsModel
	helpModel     helpModel
	error         error
}

func NewModel(client *jellyfin.Client, cfg config.Config) Model {
	return Model{
		client:        client,
		config:        cfg,
		state:         "login",
		loginModel:    newLoginModel(),
		browseModel:   newBrowseModel(),
		detailModel:   newDetailModel(),
		searchModel:   newSearchModel(),
		playlistModel: newPlaylistModel(),
		settingsModel: newSettingsModel(),
		helpModel:     newHelpModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "h":
			m.state = "help"
			return m, nil
		}
	case errors.AppError:
		m.error = msg
		return m, nil
	}

	switch m.state {
	case "login":
		m.loginModel, cmd = m.loginModel.Update(msg)
	case "browse":
		m.browseModel, cmd = m.browseModel.Update(msg)
	case "detail":
		m.detailModel, cmd = m.detailModel.Update(msg)
	case "search":
		m.searchModel, cmd = m.searchModel.Update(msg)
	case "playlist":
		m.playlistModel, cmd = m.playlistModel.Update(msg)
	case "settings":
		m.settingsModel, cmd = m.settingsModel.Update(msg)
	case "help":
		m.helpModel, cmd = m.helpModel.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	if m.error != nil {
		return errors.ErrorStyle.Render(m.error.Error())
	}

	switch m.state {
	case "login":
		return m.loginModel.View()
	case "browse":
		return m.browseModel.View()
	case "detail":
		return m.detailModel.View()
	case "search":
		return m.searchModel.View()
	case "playlist":
		return m.playlistModel.View()
	case "settings":
		return m.settingsModel.View()
	case "help":
		return m.helpModel.View()
	default:
		return "Unknown state"
	}
}
