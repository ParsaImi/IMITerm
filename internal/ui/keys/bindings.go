package keys

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
)

// Match checks if a key event matches a binding.
func Match[K fmt.Stringer](k K, b key.Binding) bool {
	return key.Matches(k, b)
}

type GroupListKeys struct {
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	Add    key.Binding
	Edit   key.Binding
	Delete key.Binding
	Sync   key.Binding
	Quit   key.Binding
}

var GroupList = GroupListKeys{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open"),
	),
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Sync: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "sync"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type HostListKeys struct {
	Back key.Binding
}

var HostList = HostListKeys{
	Back: key.NewBinding(
		key.WithKeys("esc", "b", "backspace"),
		key.WithHelp("b/esc", "back"),
	),
}

type FormKeys struct {
	Cancel key.Binding
}

var Form = FormKeys{
	Cancel: key.NewBinding(
		key.WithKeys("esc", "ctrl+b"),
		key.WithHelp("esc", "cancel"),
	),
}
