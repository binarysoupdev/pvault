package cmd

import (
	"fmt"
	"pvault/config"
	"pvault/errors"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

type ConfigCommand struct {
	command.FlagCommandBase

	configLoader config.Loader[config.Config]
	config       config.Config
}

func NewConfigCommand(loader config.Loader[config.Config]) *ConfigCommand {
	return &ConfigCommand{
		FlagCommandBase: command.NewFlagCommandBase("config", "display configuration"),
		configLoader:    loader,
	}
}

func (cmd *ConfigCommand) initialize() error {
	var err error

	cmd.config, err = cmd.configLoader.Load()
	if err != nil {
		return errors.Chain(err, "error loading config")
	}

	return nil
}

func (cmd ConfigCommand) Run(args []string) error {
	err := cmd.initialize()
	if err != nil {
		return err
	}

	style.BoldInfo.Printf("[=] Loaded from %s\n", cmd.configLoader.GetName())

	style.Bold.Printf("Version: %s\n", cmd.config.Version)
	fmt.Println("---")

	pathStyle := style.New(style.MAGENTA)
	pathStyle.Printf("Vault Path: %s\n", cmd.config.VaultPath)
	pathStyle.Printf("Output Path: %s\n", cmd.config.OutputPath)

	return nil
}
