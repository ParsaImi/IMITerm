# IMITERM

A TUI SSH connection manager for the terminal. Save your servers, organize them into groups, and connect with a single keypress. Built as a [Termius](https://termius.com/) replacement for terminal-native workflows.

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-blue)

## Features

- **Server groups** — organize hosts into named, color-coded groups
- **One-keypress connect** — select a host and SSH fires instantly (process replacement via `exec`)
- **Search** — press `/` to fuzzy-search across all hosts by name, hostname, user, group, or tags
- **Connection history** — `imiterm -h` shows your last 5 connections for quick reconnect
- **Git sync** — auto-pull on startup, auto-commit+push on every save. Sync your config across machines via a private GitHub repo
- **tmux integration** — sets the tmux window name to the host you're connected to
- **Password & key auth** — supports SSH keys, passwords (via `sshpass`), and ssh-agent

## Install

### Download binary (no Go required)

Grab the latest binary from [Releases](https://github.com/ParsaImi/imiterm/releases):

```bash
# Linux (amd64)
curl -L https://github.com/ParsaImi/imiterm/releases/latest/download/imiterm-linux-amd64 -o imiterm
chmod +x imiterm
sudo mv imiterm /usr/local/bin/

# macOS (Apple Silicon)
curl -L https://github.com/ParsaImi/imiterm/releases/latest/download/imiterm-darwin-arm64 -o imiterm
chmod +x imiterm
sudo mv imiterm /usr/local/bin/
```

### From source (requires Go 1.22+)

```bash
go install github.com/ParsaImi/imiterm/cmd/imiterm@latest
```

### System dependencies

```bash
# Required for password-authenticated hosts
sudo apt install sshpass    # Debian/Ubuntu
# or: brew install sshpass  # macOS (via homebrew)
```

## Usage

```bash
# Launch the TUI
imiterm

# Show last 5 connections (quick reconnect)
imiterm -h
```

### Key bindings

**Group list (main screen)**

| Key | Action |
|-----|--------|
| `↑/k` `↓/j` | Navigate |
| `Enter` | Open group |
| `a` | Add group |
| `e` | Edit group |
| `d` | Delete group |
| `/` | Search all hosts |
| `s` | Git sync (pull + push) |
| `q` | Quit |

**Host list**

| Key | Action |
|-----|--------|
| `Enter` | Connect (SSH) |
| `a` | Add host |
| `e` | Edit host |
| `d` | Delete host |
| `b` / `Esc` | Back (`b` is instant) |

**Forms**

| Key | Action |
|-----|--------|
| `Tab` / `↓` | Next field |
| `Shift+Tab` / `↑` | Previous field |
| `Enter` / `Ctrl+S` | Save |
| `Esc` / `Ctrl+B` | Cancel |

## Configuration

Config is stored at `~/.config/imiterm/config.toml`. Example:

```toml
[meta]
  version = 1
  git_remote = ""
  git_auto_pull = true
  git_auto_push = true

[[groups]]
  name = "Production"
  description = "Live servers"
  color = "#FF5555"

  [[groups.hosts]]
    name = "web-01"
    hostname = "192.168.1.10"
    user = "ubuntu"
    port = 22
    key_path = "~/.ssh/id_rsa"
    password = ""
    tags = ["nginx"]
```

You can edit this file directly or use the TUI forms.

### Git sync setup

Sync your config across machines via a private GitHub repo:

```bash
# 1. Create a private repo on GitHub (e.g., imiterm-config)

# 2. Initialize the config directory as a git repo
cd ~/.config/imiterm
git init
git add config.toml
git commit -m "initial config"
git remote add origin git@github.com:YOUR_USER/imiterm-config.git
git push -u origin master

# 3. Enable auto-sync in config.toml
# Set git_auto_pull = true and git_auto_push = true
```

On a new machine:

```bash
git clone git@github.com:YOUR_USER/imiterm-config.git ~/.config/imiterm
```

## Authentication

| Method | Config | How it works |
|--------|--------|-------------|
| SSH key | `key_path = "~/.ssh/id_rsa"` | `ssh -i <key> user@host` |
| Password | `password = "secret"` | `sshpass -p <pass> ssh user@host` |
| SSH agent | Both empty | `ssh user@host` (agent handles it) |

Key path takes priority if both `key_path` and `password` are set.

## License

MIT
