package model

// AuthKind tells the SSH executor which auth strategy to use.
type AuthKind int

const (
	AuthKey   AuthKind = iota // key_path is set → use ssh -i
	AuthPass                  // password is set, no key → use sshpass
	AuthAgent                 // both empty → rely on ssh-agent / default key
)

// Host is a single SSH target.
type Host struct {
	Name     string   `toml:"name"`
	Hostname string   `toml:"hostname"`
	User     string   `toml:"user"`
	Port     int      `toml:"port"`
	KeyPath  string   `toml:"key_path"`
	Password string   `toml:"password"`
	Tags     []string `toml:"tags"`
}

// AuthMethod returns which authentication strategy applies to this host.
func (h Host) AuthMethod() AuthKind {
	if h.KeyPath != "" {
		return AuthKey
	}
	if h.Password != "" {
		return AuthPass
	}
	return AuthAgent
}

// Group is a named collection of hosts.
type Group struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Color       string `toml:"color"`
	Hosts       []Host `toml:"hosts"`
}

// Meta holds global app settings stored in the config file.
type Meta struct {
	Version     int    `toml:"version"`
	GitRemote   string `toml:"git_remote"`
	GitAutoPull bool   `toml:"git_auto_pull"`
	GitAutoPush bool   `toml:"git_auto_push"`
}

// Config is the root document — what gets read from config.toml.
type Config struct {
	Meta   Meta    `toml:"meta"`
	Groups []Group `toml:"groups"`
}
