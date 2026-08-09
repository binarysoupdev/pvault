package cmds

import (
	"fmt"
	"os"
	"path/filepath"
	"pvault/app/config"
	config_flow "pvault/app/flow/config"
	vault_flow "pvault/app/flow/vault"
	"pvault/build"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/got-style/style"
)

type ConfigCommand struct {
	command.CommandBase
	command.FlagCommand
	command.ConfigCommand[config.Config]
}

func NewConfigCommand(configLoader json.Loader[config.Config]) *ConfigCommand {
	return &ConfigCommand{
		CommandBase:   command.NewCommandBase("config", "Configure the application."),
		ConfigCommand: command.NewConfigCommand(configLoader),
	}
}

func (cmd *ConfigCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	usage := cmd.Flags.Usage

	cmd.Flags.Usage = func() {
		usage()
		fmt.Println()
		cmd.printDescription()
	}

	return nil
}

func (cmd ConfigCommand) Run(args []string) error {
	new := cmd.Flags.Bool("new", false, "generate a new config file")
	cmd.ParseFlags(args)

	if *new {
		return cmd.generateNew()
	}

	err := config_flow.LoadConfig(&cmd.ConfigCommand)
	if err != nil {
		return err
	}

	return cmd.validate()
}

func (cmd ConfigCommand) printDescription() {
	style.New(style.BOLD, style.UNDERLINE).Printf("Current Version: %d\n", config.VERSION)
	fmt.Println("vault_path: the vault directory")
	fmt.Println("backup_path: directory where backups are saved")
	fmt.Println("output_path: path where unlocked files are saved")
}

func (cmd ConfigCommand) generateNew() error {
	_, err := os.Stat(cmd.ConfigLoader.Path)
	if err == nil {
		return errors.Format("config file \"%s\" already exists", cmd.ConfigLoader.Path)
	}

	base := build.DataPath()

	config := config.Config{
		Version:    config.VERSION,
		VaultPath:  filepath.Join(base, "vault"),
		BackupPath: filepath.Join(base, "backups"),
		OutputPath: ".",
	}

	err = os.MkdirAll(filepath.Dir(cmd.ConfigLoader.Path), 0755)
	if err != nil {
		return errors.Chain(err, "error creating config directory")
	}

	err = json.MarshalFilePretty(config, cmd.ConfigLoader.Path, "    ")
	if err != nil {
		return errors.Chain(err, "error saving config file")
	}

	style.BoldCreate.Printf("[+] Created New Config: %s\n", cmd.ConfigLoader.Path)
	return nil
}

func (cmd ConfigCommand) validate() error {
	style.BoldInfo.Printf("[=] Loaded from \"%s\"\n", cmd.ConfigLoader.Path)

	cmd.validateVaultPath()
	cmd.validateBackupPath()
	cmd.validateOutputPath()

	return nil
}

func (cmd ConfigCommand) validateVaultPath() {
	fmt.Printf("%s \"%s\" ", style.Bold.Sprint("Vault Path:"), cmd.Config.VaultPath)

	v, err := vault_flow.OpenCurrentVault(cmd.Config.VaultPath)
	if err != nil {
		style.Error.Printf("-> %s\n", err)
	} else {
		style.Success.Printf("-> verified (\"%s\"@v%d)\n", v.Meta.Nickname, v.GetVersion())
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
