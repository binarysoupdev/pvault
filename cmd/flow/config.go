package flow

import (
	"pvault/config"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/go-commando/json"
)

func LoadConfig(loader json.Loader[config.Config]) (config.Config, error) {
	err := loader.ValidatePath()
	if err != nil {
		return config.Config{}, errors.Chain(err, "invalid config path (run \"config -new\" to generate)")
	}

	cfg, err := loader.Load()
	if err != nil {
		return config.Config{}, errors.Chain(err, "error loading config")
	}

	if config.IsUnsupported(cfg.Version) {
		return cfg, errors.Format("unsupported version \"%d\"", cfg.Version)
	}

	if config.IsOutOfDate(cfg.Version) {
		return cfg, errors.Format("config version [%d] out-of-date (run \"config -upgrade\" to repair)", cfg.Version)
	}

	return cfg, nil
}
