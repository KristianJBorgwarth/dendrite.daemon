package persistence

import (
	"os"
	"path/filepath"
)

func getIndexDBPath(vaultPath string) (string, error) {
	dir := os.Getenv("XDG_DATA_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(home, ".local", "share")
	}

	return filepath.Join(dir, "dendrite", vaultPath), nil
}
