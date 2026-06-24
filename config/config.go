package config

import (
	"os"
	"pvault/errors"
)

type Config struct {
	Version    string `json:"version"`
	VaultPath  string `json:"vault_path"`
	OutputPath string `json:"output_path"`
}

func (c Config) ValidateOutputPath() error {
	stat, err := os.Stat(c.OutputPath)
	if err != nil {
		return errors.Chain(err, "error loading file stats")
	}

	if !stat.IsDir() {
		return errors.New("path not a directory")
	}

	return nil
}
