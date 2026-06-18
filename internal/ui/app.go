package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ParsaImi/imiterm/internal/config"
	"github.com/ParsaImi/imiterm/internal/git"
	"github.com/ParsaImi/imiterm/internal/model"
	"github.com/ParsaImi/imiterm/internal/ui/screens"
	"github.com/ParsaImi/imiterm/internal/ui/styles"
)

type screen int

const (
	screenGroupList screen = iota
	screenHostList
	screenGroupForm
	screenHostForm
	screenConfirm
	screenSearch
)

type deleteTarget int

const (
	deleteGroup deleteTarget = iota
	deleteHost
)

// syncResultMsg is returned by the background git sync command.
type syncResultMsg struct {
	err    error
	newCfg *model.Config
}

type App struct {
	cfg          *model.Config
	activeScreen screen
	prevScreen   screen

	groupList screens.GroupList
	hostList  screens.HostList
	groupForm screens.GroupForm
	hostForm  screens.HostForm
	confirm   screens.Confirm
	search    screens.Search

	activeGroupIdx int

	delTarget deleteTarget
	delIdx    int

	SelectedHost *model.Host
	status       string
	syncing      bool
}

func NewApp(cfg *model.Config) App {
	return App{
		cfg:          cfg,
		activeScreen: screenGroupList,
		groupList:    screens.NewGroupList(cfg.Groups),
	}
}

func (a App) Init() tea.Cmd { return nil }

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// --- Git sync finished (background) ---
	case syncResultMsg:
		a.syncing = false
		if msg.err != nil {
			a.status = "sync failed — saved locally"
		} else {
			a.status = "synced"
			if msg.newCfg != nil {
				a.cfg = msg.newCfg
				a.groupList = screens.NewGroupList(a.cfg.Groups)
			}
		}
		return a, nil

	case screens.GroupSelectedMsg:
		group := a.cfg.Groups[msg.Index]
		a.activeGroupIdx = msg.Index
		a.hostList = screens.NewHostList(group.Name, group.Hosts)
		a.activeScreen = screenHostList
		return a, nil

	case screens.HostSelectedMsg:
		host := msg.Host
		a.SelectedHost = &host
		return a, tea.Quit

	case screens.BackMsg:
		switch a.activeScreen {
		case screenHostList, screenSearch:
			a.activeScreen = screenGroupList
		case screenGroupForm, screenHostForm:
			if a.prevScreen == screenHostList {
				a.activeScreen = screenHostList
			} else {
				a.activeScreen = screenGroupList
			}
		case screenConfirm:
			a.activeScreen = a.prevScreen
		}
		return a, nil

	case screens.GroupSavedMsg:
		if msg.IsNew {
			a.cfg.Groups = append(a.cfg.Groups, msg.Group)
		} else {
			msg.Group.Hosts = a.cfg.Groups[msg.EditIdx].Hosts
			a.cfg.Groups[msg.EditIdx] = msg.Group
		}
		cmd := a.save()
		a.groupList = screens.NewGroupList(a.cfg.Groups)
		a.activeScreen = screenGroupList
		return a, cmd

	case screens.HostSavedMsg:
		g := &a.cfg.Groups[a.activeGroupIdx]
		if msg.IsNew {
			g.Hosts = append(g.Hosts, msg.Host)
		} else {
			g.Hosts[msg.EditIdx] = msg.Host
		}
		cmd := a.save()
		a.hostList = screens.NewHostList(g.Name, g.Hosts)
		a.activeScreen = screenHostList
		return a, cmd

	case screens.ConfirmResult:
		if msg.Confirmed {
			switch a.delTarget {
			case deleteGroup:
				a.cfg.Groups = append(a.cfg.Groups[:a.delIdx], a.cfg.Groups[a.delIdx+1:]...)
				cmd := a.save()
				a.groupList = screens.NewGroupList(a.cfg.Groups)
				a.activeScreen = screenGroupList
				return a, cmd

			case deleteHost:
				g := &a.cfg.Groups[a.activeGroupIdx]
				g.Hosts = append(g.Hosts[:a.delIdx], g.Hosts[a.delIdx+1:]...)
				cmd := a.save()
				a.hostList = screens.NewHostList(g.Name, g.Hosts)
				a.activeScreen = screenHostList
				return a, cmd
			}
		} else {
			a.activeScreen = a.prevScreen
		}
		return a, nil
	}

	switch a.activeScreen {
	case screenGroupList:
		return a.updateGroupList(msg)
	case screenHostList:
		return a.updateHostList(msg)
	case screenGroupForm:
		var cmd tea.Cmd
		a.groupForm, cmd = a.groupForm.Update(msg)
		return a, cmd
	case screenHostForm:
		var cmd tea.Cmd
		a.hostForm, cmd = a.hostForm.Update(msg)
		return a, cmd
	case screenConfirm:
		var cmd tea.Cmd
		a.confirm, cmd = a.confirm.Update(msg)
		return a, cmd
	case screenSearch:
		var cmd tea.Cmd
		a.search, cmd = a.search.Update(msg)
		return a, cmd
	}

	return a, nil
}

func (a App) updateGroupList(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		cursor := a.groupList.Cursor()

		switch msg.String() {
		case "a":
			a.groupForm = screens.NewGroupForm()
			a.prevScreen = screenGroupList
			a.activeScreen = screenGroupForm
			return a, a.groupForm.Init()

		case "e":
			if len(a.cfg.Groups) > 0 {
				a.groupForm = screens.NewGroupFormEdit(a.cfg.Groups[cursor], cursor)
				a.prevScreen = screenGroupList
				a.activeScreen = screenGroupForm
				return a, a.groupForm.Init()
			}

		case "d":
			if len(a.cfg.Groups) > 0 {
				name := a.cfg.Groups[cursor].Name
				a.confirm = screens.NewConfirm("Delete group '" + name + "'?")
				a.delTarget = deleteGroup
				a.delIdx = cursor
				a.prevScreen = screenGroupList
				a.activeScreen = screenConfirm
				return a, nil
			}

		case "s":
			cmd := a.syncNow()
			return a, cmd

		case "/":
			a.search = screens.NewSearch(a.cfg.Groups)
			a.activeScreen = screenSearch
			return a, a.search.Init()
		}
	}

	var cmd tea.Cmd
	a.groupList, cmd = a.groupList.Update(msg)
	return a, cmd
}

func (a App) updateHostList(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		cursor := a.hostList.Cursor()
		g := a.cfg.Groups[a.activeGroupIdx]

		switch msg.String() {
		case "a":
			a.hostForm = screens.NewHostForm()
			a.prevScreen = screenHostList
			a.activeScreen = screenHostForm
			return a, a.hostForm.Init()

		case "e":
			if len(g.Hosts) > 0 {
				a.hostForm = screens.NewHostFormEdit(g.Hosts[cursor], cursor)
				a.prevScreen = screenHostList
				a.activeScreen = screenHostForm
				return a, a.hostForm.Init()
			}

		case "d":
			if len(g.Hosts) > 0 {
				name := g.Hosts[cursor].Name
				a.confirm = screens.NewConfirm("Delete host '" + name + "'?")
				a.delTarget = deleteHost
				a.delIdx = cursor
				a.prevScreen = screenHostList
				a.activeScreen = screenConfirm
				return a, nil
			}
		}
	}

	var cmd tea.Cmd
	a.hostList, cmd = a.hostList.Update(msg)
	return a, cmd
}

// save writes config to disk immediately, then kicks off git sync in the background.
func (a *App) save() tea.Cmd {
	config.Save(a.cfg)

	if a.cfg.Meta.GitAutoPush && git.IsRepo(config.Dir()) {
		a.syncing = true
		a.status = "syncing..."
		dir := config.Dir()
		return func() tea.Msg {
			err := git.CommitAndPush(dir)
			return syncResultMsg{err: err}
		}
	}
	return nil
}

// syncNow does a full pull+push cycle in the background.
func (a *App) syncNow() tea.Cmd {
	dir := config.Dir()
	if !git.IsRepo(dir) {
		a.status = "not a git repo — run: cd ~/.config/imiterm && git init"
		return nil
	}

	a.syncing = true
	a.status = "syncing..."
	return func() tea.Msg {
		if err := git.Pull(dir); err != nil {
			return syncResultMsg{err: err}
		}

		newCfg, _ := config.Load()

		err := git.CommitAndPush(dir)
		return syncResultMsg{err: err, newCfg: newCfg}
	}
}

func (a App) View() string {
	var view string

	switch a.activeScreen {
	case screenGroupList:
		view = a.groupList.View()
	case screenHostList:
		view = a.hostList.View()
	case screenGroupForm:
		view = a.groupForm.View()
	case screenHostForm:
		view = a.hostForm.View()
	case screenConfirm:
		view = a.confirm.View()
	case screenSearch:
		view = a.search.View()
	}

	if a.status != "" {
		if a.syncing {
			view += styles.StatusOk.Render("  ⟳ "+a.status) + "\n"
		} else if strings.Contains(a.status, "failed") || strings.Contains(a.status, "not a git") {
			view += styles.StatusErr.Render("  ✗ "+a.status) + "\n"
		} else {
			view += styles.StatusOk.Render("  ✓ "+a.status) + "\n"
		}
	}

	return view
}
