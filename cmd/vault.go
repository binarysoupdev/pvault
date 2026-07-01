package cmd

import (
	"fmt"
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/errors"
	"pvault/vault"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

type VaultCommand struct {
	command.FlagCommandBase
	config.Loader[config.Config]
}

func NewVaultCommand(loader config.Loader[config.Config]) *VaultCommand {
	return &VaultCommand{
		FlagCommandBase: command.NewFlagCommandBase("vault", "configure the vault"),
		Loader:          loader,
	}
}

func (cmd *VaultCommand) Initialize() error {
	_ = cmd.FlagCommandBase.Initialize()
	return flow.LoadConfig(&cmd.Loader)
}

func (cmd VaultCommand) Run(args []string) error {
	init := cmd.Flags.Bool("init", false, "initialize the vault")
	upgrade := cmd.Flags.Bool("upgrade", false, "upgrade the vault if it's out-of-date")
	cmd.Flags.Parse(args)

	if *init {
		return cmd.initialize()
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

func (cmd VaultCommand) upgrade() error {
	v, err := vault.Open(cmd.Config.VaultPath)
	if err != nil {
		return errors.Chain(err, "error opening vault")
	}

	if !v.IsOutOfDate() {
		return errors.New("vault is up-to-date")
	}
	oldVersion := v.Version()

	backup, err := v.Backup(fmt.Sprintf("version_%d", v.Version()))
	if err != nil {
		return errors.Chain(err, "error backing vault")
	}

	style.BoldCreate.Printf("[+] Created Backup \"%s\"\n", backup)

	err = v.Upgrade()
	if err != nil {
		return errors.Chain(err, "error upgrading vault")
	}

	style.BoldCreate.Printf("[+] Vault Upgraded (@v%d -> @v%d)\n", oldVersion, v.Version())
	return nil
}

func (cmd VaultCommand) validate() error {
	v, err := flow.OpenVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	style.BoldInfo.Printf("[=] Vault verified at \"%s\" (@v%d)\n", v.Path, v.Version())
	style.Info.Printf("[%d] records found\n", len(v.Index))

	return nil
}
