package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Pin represents a pinned Claude Code session.
type Pin struct {
	ID          string    `json:"id"`
	Project     string    `json:"project"`
	Description string    `json:"description"`
	PinnedAt    time.Time `json:"pinned_at"`
}

// Store holds all pinned sessions.
type Store struct {
	Version  int   `json:"version"`
	Sessions []Pin `json:"sessions"`
}

func configPath() (string, error) {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(home, ".config")
	}
	return filepath.Join(dir, "pml", "pins.json"), nil
}

// Load reads the store from disk. Returns an empty store if the file doesn't exist.
func Load() (*Store, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Store{Version: 1}, nil
	}
	if err != nil {
		return nil, err
	}
	var s Store
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// Save writes the store to disk atomically.
func (s *Store) Save() error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// Upsert inserts a pin or replaces an existing one with the same ID.
func (s *Store) Upsert(p Pin) {
	for i, existing := range s.Sessions {
		if existing.ID == p.ID {
			s.Sessions[i] = p
			return
		}
	}
	s.Sessions = append(s.Sessions, p)
}

// Remove deletes a pin by ID. Returns false if not found.
func (s *Store) Remove(id string) bool {
	for i, p := range s.Sessions {
		if p.ID == id {
			s.Sessions = append(s.Sessions[:i], s.Sessions[i+1:]...)
			return true
		}
	}
	return false
}
