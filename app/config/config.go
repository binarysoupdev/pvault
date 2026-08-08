package config

const VERSION = 1

type Config struct {
	Version    int    `json:"version"`
	VaultPath  string `json:"vault_path"`
	BackupPath string `json:"backup_path"`
	OutputPath string `json:"output_path"`
}
