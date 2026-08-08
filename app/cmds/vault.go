package cmds

import (
	"path/filepath"
	"pvault/app/config"
	config_flow "pvault/app/flow/config"
	vault_flow "pvault/app/flow/vault"
	"pvault/app/vault"
	"time"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/got-style/style"
)

type VaultCommand struct {
	command.CommandBase
	command.FlagCommand
	command.ConfigCommand[config.Config]
}

func NewVaultCommand(configLoader json.Loader[config.Config]) *VaultCommand {
	return &VaultCommand{
		CommandBase:   command.NewCommandBase("vault", "configure the vault"),
		ConfigCommand: command.NewConfigCommand(configLoader),
	}
}

func (cmd *VaultCommand) Initialize() error {
	err := config_flow.LoadConfig(&cmd.ConfigCommand)
	if err != nil {
		return err
	}

	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd VaultCommand) Run(args []string) error {
	init := cmd.Flags.Bool("init", false, "initialize the vault")
	nickname := cmd.Flags.String("nickname", "", "set the vault's nickname")
	backup := cmd.Flags.Bool("backup", false, "backup the vault to the backup directory")
	upgrade := cmd.Flags.Bool("upgrade", false, "upgrade the vault if it's out-of-date")
	cmd.ParseFlags(args)

	if *init {
		return cmd.initialize(*nickname)
	}
	if *nickname != "" {
		return cmd.setNickname(*nickname)
	}
	if *backup {
		return cmd.backup()
	}
	if *upgrade {
		return cmd.upgrade()
	}
	return cmd.validate()
}

func (cmd VaultCommand) initialize(nickname string) error {
	if nickname == "" {
		nickname = filepath.Base(cmd.Config.VaultPath)
	}

	_, err := vault.InitializeNew(cmd.Config.VaultPath, nickname)
	if err != nil {
		return errors.Chain(err, "error initializing new vault")
	}

	style.BoldCreate.Printf("[+] New Vault \"%s\" Initialized: %s\n", nickname, cmd.Config.VaultPath)
	return nil
}

func (cmd VaultCommand) setNickname(nickname string) error {
	v, err := vault_flow.OpenLegacyVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	v.Meta.Nickname = nickname

	err = v.SaveMetadata()
	if err != nil {
		return err
	}

	style.Create.Printf("[+] Set Nickname: %s\n", v.Meta.Nickname)
	return nil
}

func (cmd VaultCommand) backup() error {
	v, err := vault_flow.OpenLegacyVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	return vault_flow.BackupVault(v, cmd.Config)
}

func (cmd VaultCommand) upgrade() error {
	v, err := vault_flow.OpenLegacyVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	if !v.IsOutOfDate() {
		return errors.New("vault is up-to-date")
	}
	oldVersion := v.GetVersion()

	err = vault_flow.BackupVault(v, cmd.Config)
	if err != nil {
		return err
	}

	err = v.Upgrade()
	if err != nil {
		return errors.Chain(err, "error upgrading vault")
	}

	style.BoldCreate.Printf("[+] Vault Upgraded (@v%d -> @v%d)\n", oldVersion, v.GetVersion())
	return nil
}

func (cmd VaultCommand) validate() error {
	v, err := vault_flow.OpenCurrentVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	style.BoldInfo.Printf("[=] Vault \"%s\" verified at \"%s\" (@v%d)\n", v.Meta.Nickname, v.Path, v.GetVersion())
	style.Info.Printf("Created on %s\n", v.Meta.CreationDate.Format(time.DateOnly))
	style.Info.Printf("[%d] records found\n", len(v.Map))

	return nil
}
