package cmd

import (
	"pvault/app/flow"
	"pvault/config"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/json"
)

type SearchCommand struct {
	command.CommandBase
	command.FlagCommand
	ConfigCommandBase
}

func NewSearchCommand(configLoader json.Loader[config.Config]) *SearchCommand {
	return &SearchCommand{
		CommandBase:       command.NewCommandBase("search", "search records in the vault"),
		ConfigCommandBase: NewConfigCommandBase(configLoader),
	}
}

func (cmd *SearchCommand) Initialize() error {
	if err := cmd.LoadConfig(); err != nil {
		return err
	}

	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd SearchCommand) Run(args []string) error {
	search := flow.NewSearchFlow(cmd.Flags)
	cmd.ParseFlags(args)

	v, err := flow.OpenVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	return search.Display(v)
}
