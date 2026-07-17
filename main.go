package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"pvault/app/cmd"
	"pvault/config"
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
		version.Display()
		return
	}

	configLoader := json.NewLoader[config.Config](configPath())

	runner := command.NewRunner(
		cmd.NewConfigCommand(configLoader),
		cmd.NewVaultCommand(configLoader),
		cmd.NewCreateCommand(configLoader),
		cmd.NewLockCommand(configLoader),
		cmd.NewUnlockCommand(configLoader),
		cmd.NewDeleteCommand(configLoader),
		cmd.NewSearchCommand(configLoader),
		cmd.NewCopyCommand(configLoader, clipboard.AtottoClipboard{}, qrcode.Skip2Renderer{}),
	)

	if *ls || len(os.Args) < 2 {
		version.Display()
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

	return filepath.Join(config.DataPath(), "config.json")
}
