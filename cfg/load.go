package cfg

import (
	"os"
	"path/filepath"
	"pvault/chain"
	"pvault/data"
)

func Load(cfg *Config, path string) error {
	exec, err := os.Executable()
	if err != nil {
		return chain.Error(err, "error locating executable")
	}
	path = filepath.Join(exec, path)

	*cfg, err = data.LoadJSON[Config](path)
	if err != nil {
		return chain.Error(err, "error loading config JSON")
	}

	return nil
}
