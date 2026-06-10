package cfg

type Config struct {
	VaultPath  string `json:"vault_path"`
	OutputPath string `json:"output_path"`
}

var Global Config

func SetGlobal(cfg Config) {
	Global = cfg
}
