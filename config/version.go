package config

const VERSION = 1

func (c Config) IsVersionUnsupported() bool {
	return c.Version < 1 || c.Version > VERSION
}

func (c Config) IsVersionOutOfDate() bool {
	return c.Version < VERSION
}
