package screens

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/ParsaImi/imiterm/internal/model"
	"github.com/ParsaImi/imiterm/internal/ui/keys"
	"github.com/ParsaImi/imiterm/internal/ui/styles"
)

type HistoryList struct {
	hosts    []model.Host
	cursor   int
	Selected *model.Host
}

func NewHistoryList(hosts []model.Host) HistoryList {
	return HistoryList{hosts: hosts}
}

func (h HistoryList) Init() tea.Cmd { return nil }

func (h HistoryList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
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
			host := h.hosts[h.cursor]
			h.Selected = &host
			return h, tea.Quit
		case keys.Match(msg, keys.GroupList.Quit):
			return h, tea.Quit
		}
	}
	return h, nil
}

func (h HistoryList) View() string {
	header := styles.Header.Render("  IMITERM  Recent Connections")

	var rows string
	for i, host := range h.hosts {
		rows += renderHostRow(host, i == h.cursor, "") + "\n"
	}

	content := styles.ContentBox.Render(rows)

	help := styles.RenderHelp("enter", "connect", "q", "quit")

	return header + "\n" + content + "\n" + help + "\n"
}
