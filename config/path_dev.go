//go:build !prod

package config

func DataPath() string {
	return "local"
}
