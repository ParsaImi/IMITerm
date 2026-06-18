package screens

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ParsaImi/imiterm/internal/model"
	"github.com/ParsaImi/imiterm/internal/ui/styles"
)

type GroupSavedMsg struct {
	Group   model.Group
	IsNew   bool
	EditIdx int
}

type presetColor struct {
	name string
	hex  string
}

var colorPresets = []presetColor{
	{"Red", "#FF5555"},
	{"Orange", "#FFB86C"},
	{"Yellow", "#F1FA8C"},
	{"Green", "#50FA7B"},
	{"Cyan", "#8BE9FD"},
	{"Blue", "#6272A4"},
	{"Purple", "#BD93F9"},
	{"Pink", "#FF79C6"},
}

const (
	gfName = iota
	gfDesc
	gfFieldCount
)

type GroupForm struct {
	inputs      []textinput.Model
	focused     int
	onColorRow  bool
	colorIdx    int
	isNew       bool
	editIdx     int
	err         string
}

func NewGroupForm() GroupForm {
	return newGroupForm(model.Group{Color: "#6272A4"}, true, -1)
}

func NewGroupFormEdit(g model.Group, idx int) GroupForm {
	return newGroupForm(g, false, idx)
}

func newGroupForm(g model.Group, isNew bool, editIdx int) GroupForm {
	inputs := make([]textinput.Model, gfFieldCount)

	inputs[gfName] = textinput.New()
	inputs[gfName].Placeholder = "Production"
	inputs[gfName].Prompt = ""
	inputs[gfName].Width = 40
	inputs[gfName].CharLimit = 256
	inputs[gfName].SetValue(g.Name)

	inputs[gfDesc] = textinput.New()
	inputs[gfDesc].Placeholder = "Live servers"
	inputs[gfDesc].Prompt = ""
	inputs[gfDesc].Width = 40
	inputs[gfDesc].CharLimit = 256
	inputs[gfDesc].SetValue(g.Description)

	inputs[gfName].Focus()

	colorIdx := 5
	for i, c := range colorPresets {
		if c.hex == g.Color {
			colorIdx = i
			break
		}
	}

	return GroupForm{
		inputs:   inputs,
		focused:  0,
		colorIdx: colorIdx,
		isNew:    isNew,
		editIdx:  editIdx,
	}
}

func (f GroupForm) Init() tea.Cmd { return textinput.Blink }

func (f GroupForm) Update(msg tea.Msg) (GroupForm, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "tab", "down":
			f.focusNext()
			return f, nil
		case "shift+tab", "up":
			f.focusPrev()
			return f, nil
		case "left", "h":
			if f.onColorRow {
				if f.colorIdx > 0 {
					f.colorIdx--
				}
				return f, nil
			}
		case "right", "l":
			if f.onColorRow {
				if f.colorIdx < len(colorPresets)-1 {
					f.colorIdx++
				}
				return f, nil
			}
		case "ctrl+s", "enter":
			if f.onColorRow {
				group, errMsg := f.toGroup()
				if errMsg != "" {
					f.err = errMsg
					return f, nil
				}
				return f, func() tea.Msg {
					return GroupSavedMsg{Group: group, IsNew: f.isNew, EditIdx: f.editIdx}
				}
			}
			f.focusNext()
			return f, nil
		case "esc", "ctrl+b":
			return f, func() tea.Msg { return BackMsg{} }
		}
	}

	if !f.onColorRow {
		var cmd tea.Cmd
		f.inputs[f.focused], cmd = f.inputs[f.focused].Update(msg)
		return f, cmd
	}

	return f, nil
}

func (f *GroupForm) focusNext() {
	if f.onColorRow {
		group, errMsg := f.toGroup()
		if errMsg != "" {
			f.err = errMsg
			return
		}
		_ = group
		return
	}

	f.inputs[f.focused].Blur()
	f.focused++

	if f.focused >= gfFieldCount {
		f.onColorRow = true
	} else {
		f.inputs[f.focused].Focus()
	}
}

func (f *GroupForm) focusPrev() {
	if f.onColorRow {
		f.onColorRow = false
		f.focused = gfFieldCount - 1
		f.inputs[f.focused].Focus()
		return
	}

	if f.focused > 0 {
		f.inputs[f.focused].Blur()
		f.focused--
		f.inputs[f.focused].Focus()
	}
}

func (f GroupForm) toGroup() (model.Group, string) {
	name := strings.TrimSpace(f.inputs[gfName].Value())
	if name == "" {
		return model.Group{}, "name is required"
	}

	return model.Group{
		Name:        name,
		Description: strings.TrimSpace(f.inputs[gfDesc].Value()),
		Color:       colorPresets[f.colorIdx].hex,
	}, ""
}

func (f GroupForm) View() string {
	title := "Add Group"
	if !f.isNew {
		title = "Edit Group"
	}
	header := styles.Header.Render("  IMITERM  " + title)

	labels := [gfFieldCount]string{"Name", "Description"}
	var rows string

	for i := range f.inputs {
		label := styles.FormLabel.Render(labels[i])
		indicator := "  "
		if i == f.focused && !f.onColorRow {
			indicator = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Render("▸ ")
		}
		rows += indicator + label + f.inputs[i].View() + "\n"
	}

	// Color picker row
	colorLabel := styles.FormLabel.Render("Color")
	colorIndicator := "  "
	if f.onColorRow {
		colorIndicator = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Render("▸ ")
	}

	var colorDots string
	for i, c := range colorPresets {
		dot := lipgloss.NewStyle().Foreground(lipgloss.Color(c.hex))
		if i == f.colorIdx {
			colorDots += dot.Bold(true).Render("[ ● " + c.name + " ]")
		} else {
			colorDots += dot.Render(" ● ")
		}
	}
	rows += colorIndicator + colorLabel + colorDots + "\n"

	if f.err != "" {
		rows += "\n" + styles.FormError.Render("  ✗ "+f.err) + "\n"
	}

	content := styles.FormBox.Render(rows)

	help := styles.RenderHelp(
		"tab/↓", "next", "shift+tab/↑", "prev",
		"←/→", "color", "enter", "save", "esc", "cancel",
	)

	return header + "\n" + content + "\n" + help + "\n"
}
