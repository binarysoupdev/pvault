package cmds

import (
	"pvault/app/config"
	config_flow "pvault/app/flow/config"
	vault_flow "pvault/app/flow/vault"
	"pvault/remote/host"

	crypt_host "github.com/binarysoupdev/cryptool/net/host"
	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

type HostCommand struct {
	command.CommandBase
	command.FlagCommand
	command.ConfigCommand[config.Config]
}

func NewHostCommand(configLoader json.Loader[config.Config]) *HostCommand {
	return &HostCommand{
		CommandBase:   command.NewCommandBase("host", "run the vault as a host"),
		ConfigCommand: command.NewConfigCommand(configLoader),
	}
}

func (cmd *HostCommand) Initialize() error {
	err := config_flow.LoadConfig(&cmd.ConfigCommand)
	if err != nil {
		return err
	}

	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd HostCommand) Run(args []string) error {
	port := cmd.Flags.String("port", "8080", "the port to run on")
	cmd.Flags.Parse(args)

	h, err := crypt_host.Open(*port)
	if err != nil {
		return errors.Chain(err, "error opening host")
	}
	defer h.Close()

	v, err := vault_flow.OpenCurrentVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	style.BoldInfo.Printf("Listening at: %s\n", h.Addr())
	h.Listen(host.NewHost(v), 1)

	return nil
}
