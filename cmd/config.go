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
	config.LoaderModule[config.Config]
}

func NewConfigCommand(loader config.Loader[config.Config]) *ConfigCommand {
	return &ConfigCommand{
		FlagCommandBase: command.NewFlagCommandBase("config", "display configuration"),
		LoaderModule:    config.NewLoaderModule(loader),
	}
}

func (cmd *ConfigCommand) Initialize() error {
	_ = cmd.FlagCommandBase.Initialize()

	err := cmd.LoadConfig()
	if err != nil {
		return errors.Chain(err, "error loading config")
	}

	return nil
}

func (cmd ConfigCommand) Run(args []string) error {
	style.BoldInfo.Printf("[=] Loaded from %s\n", cmd.ConfigLoader.GetName())

	style.Bold.Printf("Version: %s\n", cmd.Config.Version)
	fmt.Println("---")

	pathStyle := style.New(style.MAGENTA)
	pathStyle.Printf("Vault Path: %s\n", cmd.Config.VaultPath)
	pathStyle.Printf("Output Path: %s\n", cmd.Config.OutputPath)

	return nil
}
