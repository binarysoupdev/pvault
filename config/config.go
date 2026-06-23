package config

import (
	"os"
	"pvault/errors"
)

type Config struct {
	Path    string
	Version string `json:"version"`

	VaultPath  string `json:"vault_path"`
	OutputPath string `json:"output_path"`
}

func (c Config) Validate() error {
	errs := errors.Errors{}

	stat, err := os.Stat(c.OutputPath)
	if err != nil {
		errs.Append(errors.New("\"output_path\" invalid path"))
	} else if !stat.IsDir() {
		errs.Append(errors.New("\"output_path\" not a directory"))
	}

	return errs.Collapse()
}
