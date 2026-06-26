package cmd

import (
	"fmt"
	"pvault/cmd/flow"
	"pvault/config"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

type ConfigCommand struct {
	command.FlagCommandBase
	config.Loader[config.Config]
}

func NewConfigCommand(loader config.Loader[config.Config]) *ConfigCommand {
	return &ConfigCommand{
		FlagCommandBase: command.NewFlagCommandBase("config", "configure the application"),
		Loader:          loader,
	}
}

func (cmd ConfigCommand) Run(args []string) error {
	err := flow.LoadConfig(&cmd.Loader)
	if err != nil {
		return err
	}
	style.BoldInfo.Printf("[=] Loaded from %s\n", cmd.ConfigPath)

	style.Bold.Printf("Version [%d]", cmd.Config.Version)
	fmt.Println(" (current version)")
	fmt.Println("---")

	pathStyle := style.New(style.MAGENTA)
	pathStyle.Printf("Vault Path: %s\n", cmd.Config.VaultPath)
	pathStyle.Printf("Output Path: %s\n", cmd.Config.OutputPath)

	return nil
}
