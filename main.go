package main

import (
	"flag"
	"fmt"
	"os"
	"pvault/cmd"
	"pvault/config"
	"pvault/errors"
	"pvault/tools/clipboard"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

func main() {
	ls := flag.Bool("ls", false, "list all commands")
	flag.Parse()

	runner := command.NewRunner(
		cmd.NewConfigCommand(),
		cmd.NewInitCommand(),
		cmd.NewCreateCommand(),
		cmd.NewLockCommand(),
		cmd.NewUnlockCommand(),
		cmd.NewDeleteCommand(),
		cmd.NewSearchCommand(),
		cmd.NewCopyCommand(clipboard.AtottoClipboard{}),
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

func run(runner command.Runner) error {
	err := config.LoadDefault(&config.Global)
	if err != nil {
		return errors.Chain(err, "error loading global config")
	}

	return runner.RunCommand(os.Args[1], os.Args[2:])
}
