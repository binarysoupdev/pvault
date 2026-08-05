//go:build prod

package config

import (
	"os"
	"path/filepath"
)

func DataPath() string {
	path, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(path, ".pvault")
}
