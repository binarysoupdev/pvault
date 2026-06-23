package config

import (
	"os"
	"path/filepath"
	"pvault/data"
	"pvault/errors"
)

const CONFIG_FILE = "config.json"

func LoadDefault(cfg *Config) error {
	base, err := configPath()
	if err != nil {
		return errors.Chain(err, "error determining config path")
	}
	path := filepath.Join(base, CONFIG_FILE)

	*cfg, err = data.LoadJSON[Config](path)
	if err != nil {
		return errors.Chain(err, "error loading config JSON")
	}
	cfg.Path = path

	err = cfg.Validate()
	if err != nil {
		return errors.Chain(err, "error validating config")
	}

	return nil
}

func configPath() (string, error) {
	// check for ENV variable override
	val := os.Getenv("CFG_PATH")
	if val != "" {
		return val, nil
	}

	// use executable path as default
	exec, err := os.Executable()
	if err != nil {
		return "", errors.Chain(err, "error locating executable")
	}
	return filepath.Dir(exec), nil
}
