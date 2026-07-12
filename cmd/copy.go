package cmd

import (
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/tools/clipboard"
	"pvault/tools/qrcode"

	"github.com/binarysoupdev/go-commando/command"
	cfg "github.com/binarysoupdev/go-commando/config"
	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/got-style/style"
)

type CopyCommand struct {
	command.FlagCommandBase
	cfg.Loader[config.Config]

	clipboard clipboard.Clipboard
	qrcode    qrcode.Renderer
}

func NewCopyCommand(loader cfg.Loader[config.Config], clipboard clipboard.Clipboard, qrcode qrcode.Renderer) *CopyCommand {
	return &CopyCommand{
		FlagCommandBase: command.NewFlagCommandBase("copy", "copy password/username of a record"),
		Loader:          loader,
		clipboard:       clipboard,
		qrcode:          qrcode,
	}
}

func (cmd *CopyCommand) Initialize() error {
	_ = cmd.FlagCommandBase.Initialize()

	err := cmd.clipboard.CheckUnsupported()
	if err != nil {
		return errors.Chain(err, "clipboard unsupported")
	}

	return flow.LoadConfig(&cmd.Loader)
}

func (cmd CopyCommand) Run(args []string) error {
	search := flow.NewSearchFlow(cmd.Flags)
	username := cmd.Flags.Bool("username", false, "copy username instead of password")
	qr := cmd.Flags.Bool("qr", false, "render as a qrcode")
	cmd.Flags.Parse(args)

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

	if *qr {
		if *username {
			return cmd.renderToQRCode(USERNAME_FIELD, r.Username)
		} else {
			return cmd.renderToQRCode(PASSWORD_FIELD, r.Password)
		}
	} else {
		if *username {
			return cmd.copyToClipboard(USERNAME_FIELD, r.Username)
		} else {
			return cmd.copyToClipboard(PASSWORD_FIELD, r.Password)
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
