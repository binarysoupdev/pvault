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

	configLoader := &config.JSONLoader[config.Config]{
		Path: configPath(),
	}

	runner := command.NewRunner(
		cmd.NewConfigCommand(configLoader),
		cmd.NewInitCommand(),
		cmd.NewCreateCommand(),
		cmd.NewLockCommand(),
		cmd.NewUnlockCommand(),
		cmd.NewDeleteCommand(),
		cmd.NewSearchCommand(),
		cmd.NewCopyCommand(configLoader, clipboard.AtottoClipboard{}),
	)

	if *ls || len(os.Args) < 2 {
		runner.ListCommands()
		return
	}

	if err := run(runner); err != nil {
		style.BoldError.Print("ERROR: ")
		fmt.Println(err)
	}
}

func configPath() string {
	// check for ENV variable override
	val := os.Getenv("CFG_PATH")
	if val != "" {
		return val
	}

	// use executable path as default
	exec, _ := os.Executable()
	return filepath.Join(filepath.Dir(exec), "config.json")
}

func run(runner command.Runner) error {
	// err := config.LoadDefault(&config.Global)
	// if err != nil {
	// 	return errors.Chain(err, "error loading global config")
	// }

	return runner.RunCommand(os.Args[1], os.Args[2:])
}
