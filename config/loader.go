package config

import (
	"pvault/data"
	"pvault/errors"
)

type Loader[Config any] struct {
	ConfigPath string
	Config     Config
}

func NewLoader[Config any](path string) Loader[Config] {
	return Loader[Config]{
		ConfigPath: path,
	}
}

func (l *Loader[Config]) LoadConfig() error {
	var err error

	l.Config, err = data.LoadJSON[Config](l.ConfigPath)
	if err != nil {
		return errors.Chain(err, "error loading config JSON")
	}

	return err
}
