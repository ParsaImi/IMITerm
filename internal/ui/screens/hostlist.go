package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ParsaImi/imiterm/internal/model"
	"github.com/ParsaImi/imiterm/internal/ui/keys"
	"github.com/ParsaImi/imiterm/internal/ui/styles"
)

type HostSelectedMsg struct {
	Host model.Host
}

type BackMsg struct{}

type HostList struct {
	groupName string
	hosts     []model.Host
	cursor    int
}

func NewHostList(groupName string, hosts []model.Host) HostList {
	return HostList{groupName: groupName, hosts: hosts, cursor: 0}
}

func (h HostList) Cursor() int { return h.cursor }

func (h HostList) Init() tea.Cmd { return nil }

func (h HostList) Update(msg tea.Msg) (HostList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case keys.Match(msg, keys.GroupList.Up):
			if h.cursor > 0 {
				h.cursor--
			}
		case keys.Match(msg, keys.GroupList.Down):
			if h.cursor < len(h.hosts)-1 {
				h.cursor++
			}
		case keys.Match(msg, keys.GroupList.Enter):
			if len(h.hosts) > 0 {
				selected := h.hosts[h.cursor]
				return h, func() tea.Msg { return HostSelectedMsg{Host: selected} }
			}
		case keys.Match(msg, keys.HostList.Back):
			return h, func() tea.Msg { return BackMsg{} }
		case keys.Match(msg, keys.GroupList.Quit):
			return h, tea.Quit
		}
	}
	return h, nil
}

func (h HostList) View() string {
	header := styles.Header.Render("  IMITERM  " + h.groupName)

	var rows string
	if len(h.hosts) == 0 {
		rows = styles.SearchNoMatch.Render("No hosts in this group. Press a to add one.")
	}

	for i, host := range h.hosts {
		rows += renderHostRow(host, i == h.cursor, "") + "\n"
	}

	content := styles.ContentBox.Render(rows)

	help := styles.RenderHelp(
		"enter", "connect", "esc/b", "back",
		"a", "add", "e", "edit", "d", "delete", "q", "quit",
	)

	return header + "\n" + content + "\n" + help + "\n"
}

func renderHostRow(host model.Host, selected bool, groupName string) string {
	auth := "agent"
	switch host.AuthMethod() {
	case model.AuthKey:
		auth = "key"
	case model.AuthPass:
		auth = "pass"
	}

	name := styles.HostName.Render(fmt.Sprintf("%-15s", host.Name))
	addr := styles.HostAddr.Render(fmt.Sprintf("%s@%s:%d", host.User, host.Hostname, host.Port))
	badge := styles.RenderAuthBadge(auth)

	line := fmt.Sprintf("%s  %s  %s", name, addr, badge)

	if groupName != "" {
		tag := styles.GroupTag.Render(groupName)
		line += "  " + tag
	}

	if selected {
		return styles.ListItemSelected.Render("▸ " + line)
	}
	return styles.ListItem.Render("  " + line)
}
