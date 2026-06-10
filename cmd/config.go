package cmd

import (
	"fmt"
	"pvault/config"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

type ConfigCommand struct {
	command.FlagCommandBase
}

func NewConfigCommand() *ConfigCommand {
	return &ConfigCommand{
		FlagCommandBase: command.NewFlagCommandBase("config", "display configuration"),
	}
}

func (cmd ConfigCommand) Run(args []string) error {
	style.BoldInfo.Printf("[=] Loaded from %s\n", config.Global.Path)

	style.Bold.Printf("Version: %s\n", config.Global.Version)
	fmt.Println("---")

	pathStyle := style.New(style.MAGENTA)
	pathStyle.Printf("Vault Path: %s\n", config.Global.VaultPath)
	pathStyle.Printf("Output Path: %s\n", config.Global.OutputPath)

	return nil
}
