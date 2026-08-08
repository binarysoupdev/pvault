package cmds

import (
	"pvault/app/config"
	config_flow "pvault/app/flow/config"
	vault_flow "pvault/app/flow/vault"
	search_flow "pvault/app/flow/vault/search"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/json"
)

type SearchCommand struct {
	command.CommandBase
	command.FlagCommand
	command.ConfigCommand[config.Config]
}

func NewSearchCommand(configLoader json.Loader[config.Config]) *SearchCommand {
	return &SearchCommand{
		CommandBase:   command.NewCommandBase("search", "search records in the vault"),
		ConfigCommand: command.NewConfigCommand(configLoader),
	}
}

func (cmd *SearchCommand) Initialize() error {
	err := config_flow.LoadConfig(&cmd.ConfigCommand)
	if err != nil {
		return err
	}

	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd SearchCommand) Run(args []string) error {
	search := search_flow.NewSearchFlow(cmd.Flags)
	cmd.ParseFlags(args)

	v, err := vault_flow.OpenCurrentVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	return search.Display(v)
}
