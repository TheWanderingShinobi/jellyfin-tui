package ui

import (
	"fmt"
	"strings"

	"github.com/TheWanderingShinobi/jellyfin-tui/internal/errors"
	"github.com/TheWanderingShinobi/jellyfin-tui/internal/jellyfin"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type loginModel struct {
	focusIndex int
	inputs     []string
	cursorMode cursor
}

func newLoginModel() loginModel {
	return loginModel{
		inputs: make([]string, 2),
	}
}

type cursor int

const (
	cursorLine cursor = iota
	cursorBlock
)

func (m loginModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m loginModel) Update(msg tea.Msg) (loginModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, m.loginCmd
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = tea.Println("focused")
				} else {
					cmds[i] = tea.Println("blurred")
				}
			}

			return m, tea.Batch(cmds...)

		default:
			if m.focusIndex == len(m.inputs) {
				return m, nil
			}

			m.inputs[m.focusIndex] += msg.String()
			return m, nil
		}

	case cursor:
		m.cursorMode = msg
		return m, nil

	default:
		return m, nil
	}
}

func (m loginModel) View() string {
	var b strings.Builder

	for i := 0; i < len(m.inputs); i++ {
		b.WriteString(m.inputField(i))
		b.WriteRune('\n')
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n%s\n", *button)

	return b.String()
}

func (m loginModel) inputField(i int) string {
	var style lipgloss.Style
	if m.focusIndex == i {
		style = focusedStyle
	} else {
		style = blurredStyle
	}

	label := "Username"
	if i == 1 {
		label = "Password"
	}

	input := m.inputs[i]
	if i == 1 {
		input = strings.Repeat("*", len(input))
	}

	return fmt.Sprintf("%s: %s", style.Render(label), input)
}

func (m loginModel) loginCmd() tea.Msg {
	username := m.inputs[0]
	password := m.inputs[1]

	// Here you would typically call your Jellyfin client to perform the login
	// For this example, we'll just simulate a login
	if username == "admin" && password == "password" {
		return loginSuccessMsg{}
	}

	return errors.AppError{Type: "Login Error", Message: "Invalid username or password"}
}

type loginSuccessMsg struct{}
