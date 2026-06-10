package config

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
	path = filepath.Join(filepath.Dir(exec), path)

	*cfg, err = data.LoadJSON[Config](path)
	if err != nil {
		return chain.Error(err, "error loading config JSON")
	}

	err = cfg.Validate()
	if err != nil {
		return chain.Error(err, "error validating config")
	}

	return nil
}
