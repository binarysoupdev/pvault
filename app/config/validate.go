package config

import (
	"os"

	"github.com/binarysoupdev/go-commando/errors"
)

func (c Config) ValidateVersion() error {
	if c.Version < 1 || c.Version > VERSION {
		return errors.Format("unsupported config version \"%d\"", c.Version)
	}
	return nil
}

func (c Config) ValidateBackupPath() error {
	stat, err := os.Stat(c.BackupPath)
	if err != nil {
		return nil
	}

	if !stat.IsDir() {
		return errors.New("path not a directory")
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
