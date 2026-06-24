package config

import (
	"pvault/data"
	"pvault/errors"
)

type JSONLoader[Config any] struct {
	Path string
}

func (l JSONLoader[Config]) GetName() string {
	return l.Path
}

func (l JSONLoader[Config]) Load() (Config, error) {
	cfg, err := data.LoadJSON[Config](l.Path)
	if err != nil {
		return cfg, errors.Chain(err, "error loading config JSON")
	}

	return cfg, err
}
