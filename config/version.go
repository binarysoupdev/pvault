package config

const VERSION = 1

type Version int

func (v Version) IsUnsupported(currentVersion int) bool {
	return v < 1 || v > Version(currentVersion)
}

func (v Version) IsOutOfDate(currentVersion int) bool {
	return v < Version(currentVersion)
}
