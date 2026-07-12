package config

const (
	MIN_VERSION = 1
	VERSION     = 1
)

func IsUnsupported(version int) bool {
	return version < MIN_VERSION || version > VERSION
}

func IsOutOfDate(version int) bool {
	return version < VERSION
}
