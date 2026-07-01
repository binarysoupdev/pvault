package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/data"
	"pvault/errors"

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
	cmd.Flags.Parse(args)

	if *new {
		return cmd.generateNew()
	}

	err := flow.LoadConfig(&cmd.Loader)
	if err != nil {
		return err
	}

	return cmd.validate()
}

func (cmd ConfigCommand) generateNew() error {
	_, err := os.Stat(cmd.ConfigPath)
	if err == nil {
		return errors.Format("config file \"%s\" already exists", cmd.ConfigPath)
	}

	home := cmd.homePath()

	config := config.Config{
		Version:    config.VERSION,
		VaultPath:  filepath.Join(home, "vault"),
		BackupPath: filepath.Join(home, "backups"),
		OutputPath: ".",
	}

	err = data.SaveJSON(config, cmd.ConfigPath)
	if err != nil {
		return errors.Chain(err, "error saving config file")
	}

	style.BoldCreate.Printf("[+] Created New Config: %s\n", cmd.ConfigPath)
	return nil
}

func (ConfigCommand) homePath() string {
	base, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(base, ".pvault")
}

func (cmd ConfigCommand) validate() error {
	style.BoldInfo.Printf("[=] Loaded from \"%s\"\n", cmd.ConfigPath)

	cmd.validateVaultPath()
	cmd.validateBackupPath()
	cmd.validateOutputPath()

	return nil
}

func (cmd ConfigCommand) validateVaultPath() {
	fmt.Printf("%s \"%s\" ", style.Bold.Sprint("Vault Path:"), cmd.Config.VaultPath)

	v, err := flow.OpenVault(cmd.Config.VaultPath)
	if err != nil {
		style.Error.Printf("-> %s\n", err)
	} else {
		style.Success.Printf("-> verified (@v%d)\n", v.Version())
	}
}

func (cmd ConfigCommand) validateBackupPath() {
	fmt.Printf("%s \"%s\" ", style.Bold.Sprint("Backup Path:"), cmd.Config.BackupPath)

	err := cmd.Config.ValidateBackupPath()
	if err != nil {
		style.Error.Printf("-> %s\n", err.Error())
	} else {
		style.Success.Println("-> verified")
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
