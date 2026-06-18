package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/ParsaImi/imiterm/internal/config"
	"github.com/ParsaImi/imiterm/internal/model"
)

type Entry struct {
	Host      model.Host `json:"host"`
	Timestamp time.Time  `json:"timestamp"`
}

const maxEntries = 20

func filePath() string {
	return filepath.Join(config.Dir(), "history.json")
}

func Load() ([]Entry, error) {
	data, err := os.ReadFile(filePath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, nil
	}
	return entries, nil
}

// Record adds a host to the top of history.
func Record(h model.Host) {
	entries, _ := Load()

	entry := Entry{Host: h, Timestamp: time.Now()}
	entries = append([]Entry{entry}, entries...)

	if len(entries) > maxEntries {
		entries = entries[:maxEntries]
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(filePath(), data, 0o644)
}

// Recent returns the last n unique connections (by hostname+user+port).
func Recent(n int) []model.Host {
	entries, _ := Load()

	seen := make(map[string]bool)
	var hosts []model.Host

	for _, e := range entries {
		key := e.Host.Hostname + ":" + e.Host.User
		if seen[key] {
			continue
		}
		seen[key] = true
		hosts = append(hosts, e.Host)
		if len(hosts) >= n {
			break
		}
	}
	return hosts
}
