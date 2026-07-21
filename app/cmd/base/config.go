package base

import (
	"pvault/config"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/go-commando/json"
)

type ConfigCommand struct {
	command.ConfigCommand[config.Config]
}

func NewConfigCommand(loader json.Loader[config.Config]) ConfigCommand {
	return ConfigCommand{
		ConfigCommand: command.NewConfigCommand(loader),
	}
}

func (cmd *ConfigCommand) LoadConfig() error {
	err := cmd.ConfigLoader.ValidatePath()
	if err != nil {
		return errors.Chain(err, "invalid config path (run \"config -new\" to generate)")
	}

	err = cmd.ConfigCommand.LoadConfig()
	if err != nil {
		return err
	}

	if cmd.Config.Version.IsUnsupported(config.MIN_VERSION, config.VERSION) {
		return errors.Format("unsupported config version \"%d\"", cmd.Config.Version)
	}

	if cmd.Config.Version.IsOutOfDate(config.VERSION) {
		return errors.Format("config version \"%d\" out-of-date (run \"config -upgrade\" to repair)", cmd.Config.Version)
	}

	return nil
}
