package flow

import (
	"pvault/config"

	cfg "github.com/binarysoupdev/go-commando/config"
	"github.com/binarysoupdev/go-commando/errors"
)

func LoadConfig(loader *cfg.Loader[config.Config]) error {
	err := loader.ValidatePath()
	if err != nil {
		return errors.Chain(err, "invalid config path (run \"config -new\" to generate)")
	}

	err = loader.LoadVersion()
	if err != nil {
		return errors.Chain(err, "error loading config version")
	}

	if config.IsUnsupported(loader.ConfigVersion) {
		return errors.Format("unsupported version \"%d\"", loader.ConfigVersion)
	}

	if config.IsOutOfDate(loader.ConfigVersion) {
		return errors.Format("config version [%d] out-of-date (run \"config -upgrade\" to repair)", loader.ConfigVersion)
	}

	err = loader.LoadConfig()
	if err != nil {
		return errors.Chain(err, "error loading config")
	}

	return nil
}
