package cmds

import (
	"fmt"
	"pvault/tools/clipboard"
	"pvault/tools/rand"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/got-style/style"
)

const MAX_PASS_LENGTH = 256

type PasswordCommand struct {
	command.CommandBase
	command.FlagCommand

	clipboard clipboard.Clipboard
}

func NewPasswordCommand(clipboard clipboard.Clipboard) *PasswordCommand {
	return &PasswordCommand{
		CommandBase: command.NewCommandBase("pass", "generate a new random password"),
		clipboard:   clipboard,
	}
}

func (cmd *PasswordCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd PasswordCommand) Run(args []string) error {
	len := cmd.Flags.Int("len", 30, "length of the password")
	copy := cmd.Flags.Bool("copy", false, "copy the password directly to the clipboard")
	cmd.ParseFlags(args)

	if *len < 1 {
		return errors.New("\"len\" cannot be less than one")
	} else if *len > MAX_PASS_LENGTH {
		return errors.Format("\"len\" cannot be greater than %d", MAX_PASS_LENGTH)
	}

	password := make([]byte, *len)
	rand.Password(password)

	if *copy {
		return cmd.copyToClipboard(string(password))
	}

	fmt.Printf("%s\n", password)
	return nil
}

func (cmd PasswordCommand) copyToClipboard(password string) error {
	err := cmd.clipboard.CheckUnsupported()
	if err != nil {
		return errors.Chain(err, "clipboard unsupported")
	}

	err = cmd.clipboard.Write(password)
	if err != nil {
		return errors.Chain(err, "error copying to clipboard")
	}

	style.Create.Println("[+] PASSWORD copied to clipboard")
	return nil
}
