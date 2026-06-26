package flow

import (
	"pvault/config"
	"pvault/errors"
)

func LoadConfig(loader *config.Loader[config.Config]) error {
	err := loader.LoadConfig()
	if err != nil {
		return errors.Chain(err, "error loading config (run \"config -new\" to generate)")
	}

	err = loader.Config.ValidateVersion()
	if err != nil {
		return errors.Chain(err, "error validating config version")
	}

	if loader.Config.NeedsUpgrading() {
		return errors.Format("config version [%d] out-of-date (run \"config -upgrade\" to repair)", loader.Config.Version)
	}

	return nil
}
