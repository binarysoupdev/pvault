package config

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
	return err
}
