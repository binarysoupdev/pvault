package cmd

import (
	"pvault/app/flow"
	"pvault/app/vault"
	"pvault/config"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/got-style/style"
)

type VaultCommand struct {
	command.CommandBase
	command.FlagCommand
	flow.ConfigCommand
}

func NewVaultCommand(configLoader json.Loader[config.Config]) *VaultCommand {
	return &VaultCommand{
		CommandBase:   command.NewCommandBase("vault", "configure the vault"),
		ConfigCommand: flow.NewConfigCommand(configLoader),
	}
}

func (cmd *VaultCommand) Initialize() error {
	if err := cmd.LoadConfig(); err != nil {
		return err
	}

	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd VaultCommand) Run(args []string) error {
	init := cmd.Flags.Bool("init", false, "initialize the vault")
	backup := cmd.Flags.Bool("backup", false, "backup the vault to the backup directory")
	upgrade := cmd.Flags.Bool("upgrade", false, "upgrade the vault if it's out-of-date")
	cmd.ParseFlags(args)

	if *init {
		return cmd.initialize()
	} else if *backup {
		return cmd.backup()
	} else if *upgrade {
		return cmd.upgrade()
	} else {
		return cmd.validate()
	}
}

func (cmd VaultCommand) initialize() error {
	_, err := vault.InitializeNew(cmd.Config.VaultPath)
	if err != nil {
		return errors.Chain(err, "error initializing new vault")
	}

	style.BoldCreate.Printf("[+] New Vault Initialized: %s\n", cmd.Config.VaultPath)
	return nil
}

func (cmd VaultCommand) backup() error {
	v, err := vault.Load(cmd.Config.VaultPath)
	if err != nil {
		return errors.Chain(err, "error loading vault")
	}

	return flow.BackupVault(v, cmd.Config)
}

func (cmd VaultCommand) upgrade() error {
	v, err := vault.Load(cmd.Config.VaultPath)
	if err != nil {
		return errors.Chain(err, "error loading vault")
	}

	if !v.IsOutOfDate() {
		return errors.New("vault is up-to-date")
	}
	oldVersion := v.Version()

	err = flow.BackupVault(v, cmd.Config)
	if err != nil {
		return err
	}

	err = v.Upgrade()
	if err != nil {
		return errors.Chain(err, "error upgrading vault")
	}

	style.BoldCreate.Printf("[+] Vault Upgraded (@v%d -> @v%d)\n", oldVersion, v.Version())
	return nil
}

func (cmd VaultCommand) validate() error {
	v, err := flow.LoadVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	style.BoldInfo.Printf("[=] Vault verified at \"%s\" (@v%d)\n", v.Path, v.Version())
	style.Info.Printf("[%d] records found\n", len(v.Map))

	return nil
}
