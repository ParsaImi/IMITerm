package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ParsaImi/imiterm/internal/model"
)

// Load reads config.toml and returns the parsed Config.
// If the file doesn't exist, returns a default empty config.
func Load() (*model.Config, error) {
	cfg := &model.Config{
		Meta: model.Meta{Version: 1},
	}

	path := FilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	_, err := toml.DecodeFile(path, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Save writes the Config back to config.toml.
func Save(cfg *model.Config) error {
	path := FilePath()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(cfg)
}
