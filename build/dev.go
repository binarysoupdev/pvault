//go:build !prod

package build

func AppName() string {
	return "pvault (dev)"
}

func DataPath() string {
	return "local"
}
