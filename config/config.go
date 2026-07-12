package config

import (
	"os"

	"github.com/binarysoupdev/go-commando/errors"
)

type Config struct {
	Version    int    `json:"version"`
	VaultPath  string `json:"vault_path"`
	BackupPath string `json:"backup_path"`
	OutputPath string `json:"output_path"`
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
