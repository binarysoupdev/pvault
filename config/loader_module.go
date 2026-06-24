package config

import "pvault/errors"

type LoaderModule[Config any] struct {
	Config       Config
	ConfigLoader Loader[Config]
}

func NewLoaderModule[Config any](loader Loader[Config]) LoaderModule[Config] {
	return LoaderModule[Config]{
		ConfigLoader: loader,
	}
}

func (cfg *LoaderModule[Config]) LoadConfig() error {
	var err error

	cfg.Config, err = cfg.ConfigLoader.Load()
	if err != nil {
		return errors.Chain(err, "error loading config")
	}

	return nil
}
