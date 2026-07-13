package cmd

import (
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/tools/clipboard"
	"pvault/tools/qrcode"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/got-style/style"
)

type CopyCommand struct {
	command.CommandBase
	command.FlagCommand
	flow.ConfigCommand

	clipboard clipboard.Clipboard
	qrcode    qrcode.Renderer
}

func NewCopyCommand(configLoader json.Loader[config.Config], clipboard clipboard.Clipboard, qrcode qrcode.Renderer) *CopyCommand {
	return &CopyCommand{
		CommandBase:   command.NewCommandBase("copy", "copy password/username of a record"),
		ConfigCommand: flow.NewConfigCommand(configLoader),
		clipboard:     clipboard,
		qrcode:        qrcode,
	}
}

func (cmd *CopyCommand) Initialize() error {
	err := cmd.clipboard.CheckUnsupported()
	if err != nil {
		return errors.Chain(err, "clipboard unsupported")
	}

	if err := cmd.LoadConfig(); err != nil {
		return err
	}

	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd CopyCommand) Run(args []string) error {
	search := flow.NewSearchFlow(cmd.Flags)
	username := cmd.Flags.Bool("username", false, "copy username instead of password")
	qr := cmd.Flags.Bool("qr", false, "render as a qrcode")
	cmd.ParseFlags(args)

	v, err := flow.OpenVault(cmd.Config.VaultPath)
	if err != nil {
		return errors.Chain(err, "error opening vault")
	}

	name, err := search.Select(v)
	if err != nil {
		return err
	}

	r, err := flow.LoadVaultRecord(v, name)
	if err != nil {
		return err
	}

	const USERNAME_FIELD = "USERNAME"
	const PASSWORD_FIELD = "PASSWORD"
	r2 := r.Upgrade()

	if *qr {
		if *username {
			return cmd.renderToQRCode(USERNAME_FIELD, r2.Username)
		} else {
			return cmd.renderToQRCode(PASSWORD_FIELD, r2.Password)
		}
	} else {
		if *username {
			return cmd.copyToClipboard(USERNAME_FIELD, r2.Username)
		} else {
			return cmd.copyToClipboard(PASSWORD_FIELD, r2.Password)
		}
	}
}

func (cmd CopyCommand) renderToQRCode(field string, val string) error {
	style.Info.Printf("[=] %s rendered as QR-Code\n", field)

	err := cmd.qrcode.RenderToStdout(val)
	if err != nil {
		return errors.Chain(err, "error rendering qr-code")
	}

	return nil
}

func (cmd CopyCommand) copyToClipboard(field string, val string) error {
	err := cmd.clipboard.Write(val)
	if err != nil {
		return errors.Chain(err, "error copying to clipboard")
	}

	style.Info.Printf("[=] %s copied to clipboard\n", field)
	return nil
}
