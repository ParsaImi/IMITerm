package screens

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/ParsaImi/imiterm/internal/ui/styles"
)

type ConfirmResult struct {
	Confirmed bool
}

type Confirm struct {
	message string
}

func NewConfirm(message string) Confirm {
	return Confirm{message: message}
}

func (c Confirm) Init() tea.Cmd { return nil }

func (c Confirm) Update(msg tea.Msg) (Confirm, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "y", "Y":
			return c, func() tea.Msg { return ConfirmResult{Confirmed: true} }
		case "n", "N", "esc":
			return c, func() tea.Msg { return ConfirmResult{Confirmed: false} }
		}
	}
	return c, nil
}

func (c Confirm) View() string {
	header := styles.Header.Render("  IMITERM  Confirm")

	prompt := styles.ConfirmPrompt.Render(c.message)
	hint := styles.ListItemDim.Render("  Press y to confirm, n to cancel")

	box := styles.ConfirmBox.Render(prompt + "\n\n" + hint)

	return header + "\n" + box + "\n"
}
