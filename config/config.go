package config

import (
	"os"
	"pvault/errors"
)

const VERSION = 1

type Config struct {
	Version    int    `json:"version"`
	VaultPath  string `json:"vault_path"`
	OutputPath string `json:"output_path"`
}

func (c Config) NeedsUpgrading() bool {
	return c.Version < VERSION
}

func (c Config) ValidateVersion() error {
	if c.Version > VERSION {
		return errors.Format("unsupported version \"%d\"", c.Version)
	}
	return nil
}

func (c Config) ValidateOutputPath() error {
	stat, err := os.Stat(c.OutputPath)
	if err != nil {
		return errors.New("path not found/inaccessible")
	}

	if !stat.IsDir() {
		return errors.New("path not a directory")
	}

	return nil
}
