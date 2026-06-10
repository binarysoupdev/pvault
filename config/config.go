package config

import (
	"os"
	"pvault/chain"
)

type Config struct {
	VaultPath  string `json:"vault_path"`
	OutputPath string `json:"output_path"`
}

func (c Config) Validate() error {
	verrs := ValidationErrors{}

	stat, err := os.Stat(c.VaultPath)
	if err != nil {
		verrs = append(verrs, chain.New("\"vault_path\" invalid path"))
	} else if !stat.IsDir() {
		verrs = append(verrs, chain.New("\"vault_path\" not a directory"))
	}

	stat, err = os.Stat(c.OutputPath)
	if err != nil {
		verrs = append(verrs, chain.New("\"output_path\" invalid path"))
	} else if !stat.IsDir() {
		verrs = append(verrs, chain.New("\"output_path\" not a directory"))
	}

	if verrs.HasErrors() {
		return verrs
	}
	return nil
}
