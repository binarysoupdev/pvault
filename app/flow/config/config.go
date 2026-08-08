package config_flow

import (
	"pvault/app/config"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/errors"
)

func LoadConfig(cmd *command.ConfigCommand[config.Config]) error {
	err := cmd.ConfigLoader.ValidatePath()
	if err != nil {
		return errors.Chain(err, "invalid config path (run \"config -new\" to generate)")
	}

	err = cmd.LoadConfig()
	if err != nil {
		return errors.Chain(err, "error loading config")
	}

	err = cmd.Config.ValidateVersion()
	if err != nil {
		return err
	}

	return nil
}
