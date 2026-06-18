package screens

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/ParsaImi/imiterm/internal/model"
	"github.com/ParsaImi/imiterm/internal/ui/keys"
	"github.com/ParsaImi/imiterm/internal/ui/styles"
)

type searchResult struct {
	host      model.Host
	groupName string
}

type Search struct {
	input   textinput.Model
	all     []searchResult
	matched []searchResult
	cursor  int
}

func NewSearch(groups []model.Group) Search {
	input := textinput.New()
	input.Placeholder = "type to filter hosts..."
	input.Prompt = "  / "
	input.Focus()
	input.PromptStyle = styles.HelpKey
	input.TextStyle = styles.HostAddr

	var all []searchResult
	for _, g := range groups {
		for _, h := range g.Hosts {
			all = append(all, searchResult{host: h, groupName: g.Name})
		}
	}

	return Search{input: input, all: all, matched: all}
}

func (s Search) Init() tea.Cmd { return textinput.Blink }

func (s Search) Update(msg tea.Msg) (Search, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case keys.Match(msg, keys.GroupList.Quit):
			return s, tea.Quit
		case msg.String() == "esc" || msg.String() == "ctrl+b":
			return s, func() tea.Msg { return BackMsg{} }
		case keys.Match(msg, keys.GroupList.Up):
			if s.cursor > 0 {
				s.cursor--
			}
			return s, nil
		case keys.Match(msg, keys.GroupList.Down):
			if s.cursor < len(s.matched)-1 {
				s.cursor++
			}
			return s, nil
		case keys.Match(msg, keys.GroupList.Enter):
			if len(s.matched) > 0 {
				host := s.matched[s.cursor].host
				return s, func() tea.Msg { return HostSelectedMsg{Host: host} }
			}
			return s, nil
		}
	}

	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	s.filter()
	return s, cmd
}

func (s *Search) filter() {
	query := strings.ToLower(strings.TrimSpace(s.input.Value()))
	if query == "" {
		s.matched = s.all
		s.cursor = 0
		return
	}

	var matched []searchResult
	for _, r := range s.all {
		searchable := strings.ToLower(
			r.host.Name + " " + r.host.Hostname + " " + r.host.User + " " +
				r.groupName + " " + strings.Join(r.host.Tags, " "),
		)
		if strings.Contains(searchable, query) {
			matched = append(matched, r)
		}
	}

	s.matched = matched
	if s.cursor >= len(s.matched) {
		s.cursor = max(0, len(s.matched)-1)
	}
}

func (s Search) View() string {
	header := styles.Header.Render("  IMITERM  Search")

	searchBar := styles.SearchBox.Render(s.input.View())

	var rows string
	if len(s.matched) == 0 {
		rows = styles.SearchNoMatch.Render("no matches")
	} else {
		for i, r := range s.matched {
			rows += renderHostRow(r.host, i == s.cursor, r.groupName) + "\n"
		}
	}

	content := styles.ContentBox.Render(rows)

	help := styles.RenderHelp(
		"↑/k", "up", "↓/j", "down", "enter", "connect", "esc", "back",
	)

	return header + "\n" + searchBar + "\n" + content + "\n" + help + "\n"
}
