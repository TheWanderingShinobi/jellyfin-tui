package ui

import (
	"fmt"
	"strings"

	"github.com/TheWanderingShinobi/jellyfin-tui/internal/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	settingsTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4")).
				Padding(0, 1)

	settingsItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA"))

	settingsSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Background(lipgloss.Color("#FAFAFA"))
)

type settingsModel struct {
	cursor  int
	options []string
	config  *config.Config
}

func newSettingsModel(cfg *config.Config) settingsModel {
	return settingsModel{
		options: []string{"Server URL", "Default User", "Items Per Page"},
		config:  cfg,
	}
}

func (m settingsModel) Init() tea.Cmd {
	return nil
}

func (m settingsModel) Update(msg tea.Msg) (settingsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter":
			return m, m.editSetting
		case "esc", "q":
			return m, m.back
		}
	}
	return m, nil
}

func (m settingsModel) View() string {
	var b strings.Builder

	b.WriteString(settingsTitleStyle.Render("Settings"))
	b.WriteString("\n\n")

	for i, option := range m.options {
		value := m.getSettingValue(i)
		line := fmt.Sprintf("%s: %s", option, value)
		if i == m.cursor {
			b.WriteString(settingsSelectedStyle.Render("> " + line))
		} else {
			b.WriteString(settingsItemStyle.Render("  " + line))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString("Press Enter to edit a setting\n")
	b.WriteString("Press 'q' or Esc to go back")

	return b.String()
}

func (m settingsModel) getSettingValue(index int) string {
	switch index {
	case 0:
		return m.config.ServerURL
	case 1:
		return m.config.DefaultUser
	case 2:
		return fmt.Sprintf("%d", m.config.ItemsPerPage)
	default:
		return ""
	}
}

func (m settingsModel) editSetting() tea.Msg {
	return editSettingMsg{setting: m.options[m.cursor]}
}

func (m settingsModel) back() tea.Msg {
	return showBrowseMsg{}
}

type editSettingMsg struct {
	setting string
}

type settingsUpdateMsg struct {
	setting string
	value   string
}

// This would be a separate model for editing a single setting
type editSettingModel struct {
	setting string
	value   string
}

func newEditSettingModel(setting, initialValue string) editSettingModel {
	return editSettingModel{
		setting: setting,
		value:   initialValue,
	}
}

func (m editSettingModel) Init() tea.Cmd {
	return nil
}

func (m editSettingModel) Update(msg tea.Msg) (editSettingModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, m.save
		case "esc":
			return m, m.cancel
		case "backspace":
			if len(m.value) > 0 {
				m.value = m.value[:len(m.value)-1]
			}
		default:
			m.value += msg.String()
		}
	}
	return m, nil
}

func (m editSettingModel) View() string {
	return fmt.Sprintf(
		"Editing %s:\n\n%s\n\nPress Enter to save, Esc to cancel",
		m.setting,
		m.value,
	)
}

func (m editSettingModel) save() tea.Msg {
	return settingsUpdateMsg{setting: m.setting, value: m.value}
}

func (m editSettingModel) cancel() tea.Msg {
	return showSettingsMsg{}
}

type showSettingsMsg struct{}
