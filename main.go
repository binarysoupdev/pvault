package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	config_cmds "pvault/app/commands/config"
	record_cmds "pvault/app/commands/record"
	tool_cmds "pvault/app/commands/tool"
	vault_cmds "pvault/app/commands/vault"
	"pvault/app/config"
	"pvault/build"
	"pvault/tools/clipboard"
	"pvault/tools/qrcode"
	"pvault/version"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/got-style/style"
)

func main() {
	v := flag.Bool("v", false, "display version")
	ls := flag.Bool("ls", false, "list all commands")
	flag.Parse()

	if *v {
		version.Print()
		return
	}

	configLoader := json.NewLoader[config.Config](configPath())

	runner := command.NewRunner(
		tool_cmds.NewPasswordCommand(clipboard.AtottoClipboard{}),
		config_cmds.NewConfigCommand(configLoader),
		vault_cmds.NewVaultCommand(configLoader),
		vault_cmds.NewSearchCommand(configLoader),
		record_cmds.NewCreateCommand(configLoader),
		record_cmds.NewLockCommand(configLoader),
		record_cmds.NewUnlockCommand(configLoader),
		record_cmds.NewDeleteCommand(configLoader),
		record_cmds.NewCopyCommand(configLoader, clipboard.AtottoClipboard{}, qrcode.Skip2Renderer{}),
	)

	if *ls || len(os.Args) < 2 {
		version.Print()
		runner.ListCommands()
		return
	}

	if err := runner.RunCommand(os.Args[1], os.Args[2:]); err != nil {
		style.BoldError.Print("ERROR: ")
		fmt.Println(err)
	}
}

func configPath() string {
	val := os.Getenv("CONFIG")
	if val != "" {
		return val
	}

	return filepath.Join(build.DataPath(), "config.json")
}
