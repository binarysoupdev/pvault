package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/data"
	"pvault/errors"
	"pvault/vault"

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
	new := cmd.Flags.Bool("new", false, "generate a new config file")
	init := cmd.Flags.Bool("init", false, "initialize the vault")
	upgrade := cmd.Flags.Bool("upgrade", false, "upgrade the vault if it's out-of-date")
	cmd.Flags.Parse(args)

	if *new {
		return cmd.generateNew()
	}

	err := flow.LoadConfig(&cmd.Loader)
	if err != nil {
		return err
	}

	if *init {
		return cmd.initializeVault()
	} else if *upgrade {
		return cmd.upgradeVault()
	} else {
		return cmd.validate()
	}
}

func (cmd ConfigCommand) generateNew() error {
	_, err := os.Stat(cmd.ConfigPath)
	if err == nil {
		return errors.Format("config file \"%s\" already exists", cmd.ConfigPath)
	}

	config := config.Config{
		Version:    config.VERSION,
		VaultPath:  cmd.newVaultPath(),
		OutputPath: ".",
	}

	err = data.SaveJSON(config, cmd.ConfigPath)
	if err != nil {
		return errors.Chain(err, "error saving config file")
	}

	style.BoldCreate.Printf("[+] Created New Config: %s\n", cmd.ConfigPath)
	return nil
}

func (ConfigCommand) newVaultPath() string {
	base, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(base, ".pvault/vault")
}

func (cmd ConfigCommand) initializeVault() error {
	_, err := vault.InitializeNew(cmd.Config.VaultPath)
	if err != nil {
		return errors.Chain(err, "error initializing new vault")
	}

	style.BoldCreate.Printf("[+] New Vault Initialized: %s\n", cmd.Config.VaultPath)
	return nil
}

func (cmd ConfigCommand) upgradeVault() error {
	v, err := vault.Open(cmd.Config.VaultPath)
	if err != nil {
		return errors.Chain(err, "error opening vault")
	}

	if !v.IsOutOfDate() {
		return errors.New("vault is up-to-date")
	}

	backup := fmt.Sprintf("version_%d", v.Version())
	err = v.Backup(backup)
	if err != nil {
		return errors.Chain(err, "error backing vault")
	}

	style.BoldCreate.Printf("[+] Created Backup: %s\n", backup)
	return nil
}

func (cmd ConfigCommand) validate() error {
	style.BoldInfo.Printf("[=] Loaded from \"%s\"\n", cmd.ConfigPath)

	cmd.validateVaultPath()
	cmd.validateOutputPath()

	return nil
}

func (cmd ConfigCommand) validateVaultPath() {
	fmt.Printf("%s \"%s\" ", style.Bold.Sprint("Vault Path:"), cmd.Config.VaultPath)

	v, err := vault.Open(cmd.Config.VaultPath)
	if err != nil {
		style.Error.Println("-> error opening vault (run \"config -init\" to repair)")
		return
	}

	if v.IsOutOfDate() {
		style.Error.Printf("-> vault (@v%d) out-of-date (run \"config -upgrade\" to repair)\n", v.Version())
	} else {
		style.Success.Printf("-> verified (@v%d)\n", v.Version())
	}
}

func (cmd ConfigCommand) validateOutputPath() {
	fmt.Printf("%s \"%s\" ", style.Bold.Sprint("Output Path:"), cmd.Config.OutputPath)

	err := cmd.Config.ValidateOutputPath()
	if err != nil {
		style.Error.Printf("-> %s\n", err.Error())
	} else {
		style.Success.Println("-> verified")
	}
}
