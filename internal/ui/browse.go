package ui

import (
	"fmt"
	"strings"

	"github.com/TheWanderingShinobi/jellyfin-tui/internal/jellyfin"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type browseModel struct {
	items        []jellyfin.MediaItem
	cursor       int
	selected     map[int]struct{}
	client       *jellyfin.Client
	filter       string
	page         int
	itemsPerPage int
	totalItems   int
}

func newBrowseModel(client *jellyfin.Client) browseModel {
	return browseModel{
		selected:     make(map[int]struct{}),
		client:       client,
		page:         1,
		itemsPerPage: 20,
	}
}

func (m browseModel) Init() tea.Cmd {
	return m.fetchItems
}

func (m browseModel) Update(msg tea.Msg) (browseModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "n":
			if (m.page * m.itemsPerPage) < m.totalItems {
				m.page++
				return m, m.fetchItems
			}
		case "p":
			if m.page > 1 {
				m.page--
				return m, m.fetchItems
			}
		case "f":
			return m, m.showFilter
		case "s":
			return m, m.showSearch
		case "q", "esc":
			return m, m.quit
		}
	case mediaItemsMsg:
		m.items = msg.items
		m.totalItems = msg.totalItems
	}

	return m, nil
}

func (m browseModel) View() string {
	s := titleStyle.Render("Browse Media Items") + "\n\n"

	for i, item := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		style := itemStyle
		if m.cursor == i {
			style = selectedItemStyle
		}

		s += style.Render(fmt.Sprintf("%s [%s] %s\n", cursor, checked, item.Name))
	}

	s += fmt.Sprintf("\nPage %d of %d", m.page, (m.totalItems+m.itemsPerPage-1)/m.itemsPerPage)
	s += "\n\nPress 'f' to filter, 's' to search, 'n' for next page, 'p' for previous page"
	s += "\nPress 'q' to quit"

	return s
}

func (m browseModel) fetchItems() tea.Msg {
	items, total, err := m.client.GetMediaItems(m.page, m.itemsPerPage, m.filter)
	if err != nil {
		return errorMsg{err}
	}
	return mediaItemsMsg{items: items, totalItems: total}
}

func (m browseModel) showFilter() tea.Msg {
	return showFilterMsg{}
}

func (m browseModel) showSearch() tea.Msg {
	return showSearchMsg{}
}

func (m browseModel) quit() tea.Msg {
	return quitMsg{}
}

type mediaItemsMsg struct {
	items      []jellyfin.MediaItem
	totalItems int
}

type showFilterMsg struct{}
type showSearchMsg struct{}
type quitMsg struct{}
