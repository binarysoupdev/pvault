package main

// go build -tags=prod -ldflags="-s -w" -o pvault-prod

import (
	"fmt"
	"os"
	"path/filepath"
	"pvault/app/cmds"
	"pvault/app/config"
	"pvault/app/logger"
	"pvault/build"
	"pvault/tools/clipboard"
	"pvault/tools/qrcode"
	"pvault/version"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/got-style/style"
)

func main() {
	runner := buildRunner()

	if len(os.Args) < 2 {
		printDefault(runner)
		return
	}

	if err := runApp(runner); err != nil {
		style.BoldError.Print("ERROR: ")
		fmt.Println(err)
	}
}

func printDefault(runner command.Runner) {
	version.Print()
	runner.ListCommands()
}

func runApp(runner command.Runner) error {
	log, err := logger.Open(logPath())
	if err != nil {
		return errors.Chain(err, "error opening logger")
	}
	defer log.Close()

	return runner.RunCommand(os.Args[1], os.Args[2:])
}

//=================================================

func buildRunner() command.Runner {
	configLoader := json.NewLoader[config.Config](configPath())
	clipboard := clipboard.AtottoClipboard{}
	qrcode := qrcode.Skip2Renderer{}

	return command.NewRunner(
		cmds.NewConfigCommand(configLoader),
		cmds.NewVaultCommand(configLoader),
		cmds.NewPasswordCommand(clipboard),
		cmds.NewSearchCommand(configLoader),
		cmds.NewCreateCommand(configLoader),
		cmds.NewLockCommand(configLoader),
		cmds.NewUnlockCommand(configLoader),
		cmds.NewDeleteCommand(configLoader),
		cmds.NewCopyCommand(configLoader, clipboard, qrcode),
	)
}

func configPath() string {
	val := os.Getenv("CONFIG")
	if val != "" {
		return val
	}

	return filepath.Join(build.DataPath(), "config.json")
}

func logPath() string {
	val := os.Getenv("LOG")
	if val != "" {
		return val
	}

	return filepath.Join(build.DataPath(), "log.txt")
}
