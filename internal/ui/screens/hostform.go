package screens

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ParsaImi/imiterm/internal/model"
	"github.com/ParsaImi/imiterm/internal/ui/styles"
)

type HostSavedMsg struct {
	Host    model.Host
	IsNew   bool
	EditIdx int
}

const (
	fieldName = iota
	fieldHostname
	fieldUser
	fieldPort
	fieldKeyPath
	fieldPassword
	fieldTags
	fieldCount
)

var fieldLabels = [fieldCount]string{
	"Name", "Hostname", "User", "Port", "Key Path", "Password", "Tags",
}

type HostForm struct {
	inputs  []textinput.Model
	focused int
	isNew   bool
	editIdx int
	err     string
}

func NewHostForm() HostForm {
	return newHostForm(model.Host{Port: 22}, true, -1)
}

func NewHostFormEdit(h model.Host, idx int) HostForm {
	return newHostForm(h, false, idx)
}

func newHostForm(h model.Host, isNew bool, editIdx int) HostForm {
	inputs := make([]textinput.Model, fieldCount)

	placeholders := [fieldCount]string{
		"web-01", "192.168.1.10", "ubuntu", "22",
		"~/.ssh/id_rsa", "(leave empty for key/agent)", "nginx, frontend",
	}

	for i := range inputs {
		inputs[i] = textinput.New()
		inputs[i].CharLimit = 256
		inputs[i].Placeholder = placeholders[i]
		inputs[i].Prompt = ""
		inputs[i].Width = 40
	}

	inputs[fieldName].SetValue(h.Name)
	inputs[fieldHostname].SetValue(h.Hostname)
	inputs[fieldUser].SetValue(h.User)
	inputs[fieldPort].SetValue(fmt.Sprintf("%d", h.Port))
	inputs[fieldKeyPath].SetValue(h.KeyPath)
	inputs[fieldPassword].SetValue(h.Password)
	inputs[fieldTags].SetValue(strings.Join(h.Tags, ", "))

	inputs[fieldName].Focus()

	return HostForm{inputs: inputs, focused: 0, isNew: isNew, editIdx: editIdx}
}

func (f HostForm) Init() tea.Cmd { return textinput.Blink }

func (f HostForm) Update(msg tea.Msg) (HostForm, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "tab", "down":
			f.focusNext()
			return f, nil
		case "shift+tab", "up":
			f.focusPrev()
			return f, nil
		case "ctrl+s", "enter":
			host, err := f.toHost()
			if err != "" {
				f.err = err
				return f, nil
			}
			return f, func() tea.Msg {
				return HostSavedMsg{Host: host, IsNew: f.isNew, EditIdx: f.editIdx}
			}
		case "esc", "ctrl+b":
			return f, func() tea.Msg { return BackMsg{} }
		}
	}

	var cmd tea.Cmd
	f.inputs[f.focused], cmd = f.inputs[f.focused].Update(msg)
	return f, cmd
}

func (f *HostForm) focusNext() {
	f.inputs[f.focused].Blur()
	f.focused = (f.focused + 1) % fieldCount
	f.inputs[f.focused].Focus()
}

func (f *HostForm) focusPrev() {
	f.inputs[f.focused].Blur()
	f.focused = (f.focused - 1 + fieldCount) % fieldCount
	f.inputs[f.focused].Focus()
}

func (f HostForm) toHost() (model.Host, string) {
	name := strings.TrimSpace(f.inputs[fieldName].Value())
	hostname := strings.TrimSpace(f.inputs[fieldHostname].Value())
	user := strings.TrimSpace(f.inputs[fieldUser].Value())
	portStr := strings.TrimSpace(f.inputs[fieldPort].Value())
	keyPath := strings.TrimSpace(f.inputs[fieldKeyPath].Value())
	password := f.inputs[fieldPassword].Value()
	tagsRaw := strings.TrimSpace(f.inputs[fieldTags].Value())

	if name == "" {
		return model.Host{}, "name is required"
	}
	if hostname == "" {
		return model.Host{}, "hostname is required"
	}
	if user == "" {
		return model.Host{}, "user is required"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return model.Host{}, "port must be 1-65535"
	}

	var tags []string
	if tagsRaw != "" {
		for _, t := range strings.Split(tagsRaw, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, t)
			}
		}
	}

	return model.Host{
		Name: name, Hostname: hostname, User: user, Port: port,
		KeyPath: keyPath, Password: password, Tags: tags,
	}, ""
}

func (f HostForm) View() string {
	title := "Add Host"
	if !f.isNew {
		title = "Edit Host"
	}
	header := styles.Header.Render("  IMITERM  " + title)

	var rows string
	for i := range f.inputs {
		label := styles.FormLabel.Render(fieldLabels[i])
		indicator := "  "
		if i == f.focused {
			indicator = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Render("▸ ")
		}
		rows += indicator + label + f.inputs[i].View() + "\n"
	}

	if f.err != "" {
		rows += "\n" + styles.FormError.Render("  ✗ "+f.err) + "\n"
	}

	content := styles.FormBox.Render(rows)

	help := styles.RenderHelp(
		"tab/↓", "next", "shift+tab/↑", "prev",
		"enter/ctrl+s", "save", "esc", "cancel",
	)

	return header + "\n" + content + "\n" + help + "\n"
}
