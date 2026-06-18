# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
go build -o imiterm ./cmd/imiterm/     # build binary
go run ./cmd/imiterm/                   # dev run (needs a real TTY)
go install ./cmd/imiterm/               # install to $GOPATH/bin
```

No test suite yet. Verify manually by running the binary in a terminal — piped input won't work (Bubble Tea requires a TTY).

## Architecture

Elm-architecture TUI built with Bubble Tea (charmbracelet/bubbletea).

**Data flow:** `main.go` loads config → optional git pull → starts Bubble Tea → user navigates → on host select, TUI exits → `syscall.Exec` replaces process with SSH.

**Key constraint:** `syscall.Exec` must happen AFTER `p.Run()` returns, never inside a `tea.Cmd`. Bubble Tea restores the terminal on exit; exec inside Update would corrupt it.

**Package dependency graph** (arrows = imports):
```
cmd/imiterm/main.go → config, git, history, ssh, ui
ui/app.go → config, git, model, ui/screens, ui/styles
ui/screens/* → model, ui/keys, ui/styles
ssh/exec.go → model
config/ → model
git/ → (no internal deps)
history/ → config, model
model/ → (no internal deps, center of the graph)
```

**Screen routing:** `ui/app.go` holds a `screen` enum and dispatches `Update`/`View` to the active screen. Screens communicate back via typed `tea.Msg` values (e.g., `HostSelectedMsg`, `BackMsg`). Navigation messages are handled in `App.Update` before delegation.

**Git sync:** Runs as a `tea.Cmd` (background goroutine) to avoid freezing the TUI. Returns `syncResultMsg` when done.

## Config

TOML at `~/.config/imiterm/config.toml`. Structs in `internal/model/types.go`. Load/Save in `internal/config/config.go`.
