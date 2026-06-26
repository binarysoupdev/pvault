package cmd

import (
	"fmt"
	"pvault/cmd/flow"
	"pvault/config"
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
	err := flow.LoadConfig(&cmd.Loader)
	if err != nil {
		return err
	}
	style.BoldInfo.Printf("[=] Loaded from \"%s\"\n", cmd.ConfigPath)

	fmt.Printf("%s (current version)\n", style.New(style.MAGENTA, style.BOLD).Sprintf("Version [%d]", cmd.Config.Version))
	fmt.Println("---")

	cmd.validateVaultPath()
	cmd.validateOutputPath()

	return nil
}

func (cmd ConfigCommand) validateVaultPath() {
	fmt.Printf("%s \"%s\" ", style.Bold.Sprint("Vault Path:"), cmd.Config.VaultPath)

	_, err := vault.Open(cmd.Config.VaultPath)
	if err != nil {
		style.Error.Println("-> error opening vault (run \"config -init\" to repair)")
	} else {
		fmt.Printf("(vault@v%d)", vault.VERSION)
		style.Success.Println(" -> verified!")
	}
}

func (cmd ConfigCommand) validateOutputPath() {
	fmt.Printf("%s \"%s\" ", style.Bold.Sprint("Output Path:"), cmd.Config.OutputPath)

	err := cmd.Config.ValidateOutputPath()
	if err != nil {
		style.Error.Printf("-> %s\n", err.Error())
	} else {
		style.Success.Println("-> verified!")
	}
}
