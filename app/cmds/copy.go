package cmds

import (
	"fmt"
	"pvault/app/config"
	config_flow "pvault/app/flow/config"
	vault_flow "pvault/app/flow/vault"
	search_flow "pvault/app/flow/vault/search"
	"pvault/tools/clipboard"
	"pvault/tools/qrcode"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

type CopyCommand struct {
	command.CommandBase
	command.FlagCommand
	command.ConfigCommand[config.Config]

	clipboard clipboard.Clipboard
	qrcode    qrcode.Renderer
}

func NewCopyCommand(configLoader json.Loader[config.Config], clipboard clipboard.Clipboard, qrcode qrcode.Renderer) *CopyCommand {
	return &CopyCommand{
		CommandBase:   command.NewCommandBase("copy", "Copy password/username of a record"),
		ConfigCommand: command.NewConfigCommand(configLoader),
		clipboard:     clipboard,
		qrcode:        qrcode,
	}
}

func (cmd *CopyCommand) Initialize() error {
	err := cmd.clipboard.CheckUnsupported()
	if err != nil {
		return errors.Chain(err, "clipboard unsupported")
	}

	err = config_flow.LoadConfig(&cmd.ConfigCommand)
	if err != nil {
		return err
	}

	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd CopyCommand) Run(args []string) error {
	search := search_flow.NewSearchFlow(cmd.Flags)
	username := cmd.Flags.Bool("username", false, "copy username instead of password")
	qr := cmd.Flags.Bool("qr", false, "render as a qrcode")
	show := cmd.Flags.Bool("show", false, "print to console")
	cmd.ParseFlags(args)

	v, err := vault_flow.OpenCurrentVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	name, err := search.Select(v)
	if err != nil {
		return err
	}

	r, err := vault_flow.LoadRecord(v, name)
	if err != nil {
		return err
	}

	var field string
	if *username {
		style.Info.Print("[=] USERNAME")
		field = r.Username
	} else {
		style.Info.Print("[=] PASSWORD")
		field = r.Password
	}

	switch {
	case *qr:
		return cmd.renderToQRCode(field)
	case *show:
		return cmd.printField(field)
	default:
		return cmd.copyToClipboard(field)
	}
}

func (cmd CopyCommand) renderToQRCode(field string) error {
	style.Info.Println("rendered as QR-Code")

	err := cmd.qrcode.RenderToStdout(field)
	if err != nil {
		return errors.Chain(err, "error rendering qr-code")
	}

	return nil
}

func (cmd CopyCommand) printField(field string) error {
	fmt.Printf(": %s\n", field)
	return nil
}

func (cmd CopyCommand) copyToClipboard(field string) error {
	style.Info.Println("copied to clipboard")

	err := cmd.clipboard.Write(field)
	if err != nil {
		return errors.Chain(err, "error copying to clipboard")
	}

	return nil
}
