//go:build prod

package build

import (
	"os"
	"path/filepath"
)

func AppName() string {
	return "pvault"
}

func DataPath() string {
	path, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(path, ".pvault")
}
