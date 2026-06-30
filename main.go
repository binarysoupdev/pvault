package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"pvault/cmd"
	"pvault/config"
	"pvault/tools/clipboard"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

func main() {
	ls := flag.Bool("ls", false, "list all commands")
	flag.Parse()

	configLoader := config.NewLoader[config.Config](configPath())

	runner := command.NewRunner(
		cmd.NewConfigCommand(configLoader),
		cmd.NewVaultCommand(configLoader),
		cmd.NewCreateCommand(configLoader),
		cmd.NewLockCommand(configLoader),
		cmd.NewUnlockCommand(configLoader),
		cmd.NewDeleteCommand(configLoader),
		cmd.NewSearchCommand(configLoader),
		cmd.NewCopyCommand(configLoader, clipboard.AtottoClipboard{}),
	)

	if *ls || len(os.Args) < 2 {
		runner.ListCommands()
		return
	}

	if err := runner.RunCommand(os.Args[1], os.Args[2:]); err != nil {
		style.BoldError.Print("ERROR: ")
		fmt.Println(err)
	}
}

func configPath() string {
	// check for ENV variable override
	val := os.Getenv("CONFIG")
	if val != "" {
		return val
	}

	// use executable path as default
	exec, _ := os.Executable()
	return filepath.Join(filepath.Dir(exec), "config.json")
}
