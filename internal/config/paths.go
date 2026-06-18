package config

import (
	"os"
	"path/filepath"
)

// Dir returns the config directory (~/.config/imiterm/).
// Creates it if it doesn't exist.
func Dir() string {
	base, err := os.UserConfigDir()
	if err != nil {
		base = filepath.Join(os.Getenv("HOME"), ".config")
	}
	dir := filepath.Join(base, "imiterm")
	os.MkdirAll(dir, 0o755)
	return dir
}

// FilePath returns the full path to config.toml.
func FilePath() string {
	return filepath.Join(Dir(), "config.toml")
}
