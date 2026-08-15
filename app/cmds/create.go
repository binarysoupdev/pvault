package cmds

import (
	"pvault/app/config"
	config_flow "pvault/app/flow/config"
	output_flow "pvault/app/flow/output"
	vault_flow "pvault/app/flow/vault"
	record_v2 "pvault/vault/record/record/v2"

	"github.com/binarysoupdev/cryptool/rand"
	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
)

type CreateCommand struct {
	command.CommandBase
	command.FlagCommand
	command.ConfigCommand[config.Config]
}

func NewCreateCommand(configLoader json.Loader[config.Config]) *CreateCommand {
	return &CreateCommand{
		CommandBase:   command.NewCommandBase("create", "Create a new vault record"),
		ConfigCommand: command.NewConfigCommand(configLoader),
	}
}

func (cmd *CreateCommand) Initialize() error {
	err := config_flow.LoadConfig(&cmd.ConfigCommand)
	if err != nil {
		return err
	}

	err = cmd.Config.ValidateOutputPath()
	if err != nil {
		return errors.Chain(err, "error validating \"config.output_path\"")
	}

	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd CreateCommand) Run(args []string) error {
	name := cmd.Flags.String("name", "", "name of the record")
	pass := cmd.Flags.Int("pass", 0, "populate the password field with a random n-len string")
	cmd.ParseFlags(args)

	v, err := vault_flow.OpenCurrentVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	if *name == "" {
		return errors.New("\"name\" cannot be empty")
	}
	if *pass > MAX_PASS_LENGTH {
		return errors.Format("\"pass\" cannot be greater than %d", MAX_PASS_LENGTH)
	}

	r := record_v2.NewEmptyRecord(*name)

	if *pass > 0 {
		password := make([]byte, *pass)
		rand.ASCII(password)

		r.Password = string(password)
	}

	err = v.ValidateRecord(r)
	if err != nil {
		return errors.Format("name \"%s\" already exists", *name)
	}

	err = vault_flow.CreateRecord(v, r)
	if err != nil {
		return err
	}

	err = output_flow.SaveRecord(cmd.Config, r)
	if err != nil {
		return err
	}

	return nil
}
