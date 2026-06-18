package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ParsaImi/imiterm/internal/model"
	"github.com/ParsaImi/imiterm/internal/ui/keys"
	"github.com/ParsaImi/imiterm/internal/ui/styles"
)

type GroupSelectedMsg struct {
	Index int
}

type GroupList struct {
	groups []model.Group
	cursor int
	width  int
	height int
}

func NewGroupList(groups []model.Group) GroupList {
	return GroupList{groups: groups, cursor: 0}
}

func (g GroupList) Cursor() int { return g.cursor }

func (g GroupList) Init() tea.Cmd { return nil }

func (g GroupList) Update(msg tea.Msg) (GroupList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		g.width = msg.Width
		g.height = msg.Height
	case tea.KeyMsg:
		switch {
		case keys.Match(msg, keys.GroupList.Up):
			if g.cursor > 0 {
				g.cursor--
			}
		case keys.Match(msg, keys.GroupList.Down):
			if g.cursor < len(g.groups)-1 {
				g.cursor++
			}
		case keys.Match(msg, keys.GroupList.Enter):
			if len(g.groups) > 0 {
				return g, func() tea.Msg { return GroupSelectedMsg{Index: g.cursor} }
			}
		case keys.Match(msg, keys.GroupList.Quit):
			return g, tea.Quit
		}
	}
	return g, nil
}

func (g GroupList) View() string {
	header := styles.Header.Render("  IMITERM  SSH Manager")

	var rows string
	if len(g.groups) == 0 {
		rows = styles.SearchNoMatch.Render("No groups yet. Press a to add one.")
	}

	for i, group := range g.groups {
		dot := lipgloss.NewStyle().Foreground(lipgloss.Color(group.Color)).Render("●")
		name := styles.GroupName.Render(group.Name)
		count := styles.GroupCount.Render(fmt.Sprintf("%d hosts", len(group.Hosts)))

		line := fmt.Sprintf("%s  %s  %s", dot, name, count)

		if group.Description != "" {
			desc := styles.ListItemDim.Render(" — " + group.Description)
			line += desc
		}

		if i == g.cursor {
			rows += styles.ListItemSelected.Render("▸ "+line) + "\n"
		} else {
			rows += styles.ListItem.Render("  "+line) + "\n"
		}
	}

	content := styles.ContentBox.Render(rows)

	help := styles.RenderHelp(
		"↑/k", "up", "↓/j", "down", "enter", "open",
		"a", "add", "e", "edit", "d", "delete",
		"/", "search", "s", "sync", "q", "quit",
	)

	return header + "\n" + content + "\n" + help + "\n"
}
