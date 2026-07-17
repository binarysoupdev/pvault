package cmd_test

import (
	"errors"
	"os"
	vault "pvault/app/vault/local"
	v2 "pvault/app/vault/record/version2"
	"pvault/cmd"
	"pvault/config"
	"pvault/tools/clipboard"
	"pvault/tools/qrcode"
	"testing"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type CopyTestSuite struct {
	test.CommandSuite[*cmd.CopyCommand]
	Clipboard *clipboard.MockClipboard
	QRCode    *qrcode.MockRenderer

	ConfigLoader json.Loader[config.Config]
	Config       config.Config

	Vault    vault.Vault
	Record   v2.Record
	Password string
}

func TestCopyCommandSuite(t *testing.T) {
	s := CopyTestSuite{
		ConfigLoader: json.NewLoader[config.Config](file.NewPath(t, "config.json")),
		Clipboard:    clipboard.Mock(),
		QRCode:       qrcode.Mock(),
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewCopyCommand(s.ConfigLoader, s.Clipboard, s.QRCode))
	suite.Run(t, &s)
}

func (s *CopyTestSuite) SetupTest() {
	*s.Clipboard = clipboard.MockClipboard{}
	*s.QRCode = qrcode.MockRenderer{}

	s.Config = config.Config{
		Version:   config.VERSION,
		VaultPath: file.NewPath(s.T(), "vault"),
	}
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	rand := rand.New(0)
	s.Record = v2.NewEmptyRecord(rand.ASCII(15))
	s.Record.Username = rand.ASCII(10)
	s.Record.Password = rand.ASCII(30)

	s.Password = rand.ASCII(30)

	s.Vault, err = vault.InitializeNew(s.Config.VaultPath)
	s.Require().NoError(err)

	err = s.Vault.SaveRecord(s.Record, s.Password)
	s.Require().NoError(err)
}

//=====================================

func (s *CopyTestSuite) TestRunFailClipboardUnsupported() {
	//-- arrange
	s.Clipboard.Unsupported = errors.New("")

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("clipboard unsupported")
}

func (s *CopyTestSuite) TestRunFailErrorLoadingConfig() {
	//-- arrange
	err := os.Remove(s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("invalid config path")
}

func (s *CopyTestSuite) TestRunFailInvalidVaultPath() {
	//-- arrange
	s.Config.VaultPath = "invalid"
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *CopyTestSuite) TestRunFailInvalidNoResults() {
	//-- act
	s.RunCommand("-s", "no match")

	//-- assert
	s.RequireResultFail("no matches found")
}

func (s *CopyTestSuite) TestRunFailIncorrectPassword() {
	//-- arrange
	io := pipe.OpenStdio(1, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", s.Password+"x")
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error decrypting record")

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Enter PASSWORD")
}

func (s *CopyTestSuite) TestRunFailErrorCopyingToClipboard() {
	//-- arrange
	s.Clipboard.Error = errors.New("")

	io := pipe.OpenStdio(1, 3, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", s.Password)
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error copying to clipboard")

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Enter PASSWORD")
	s.Assert().Contains(io.ReadLine(), "[=] Loaded Record: "+s.Record.ID.String())
}

func (s *CopyTestSuite) TestRunPassPasswordCopied() {
	//-- arrange
	io := pipe.OpenStdio(1, 4, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", s.Password)
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultPass()

	s.Assert().Equal(s.Record.Password, s.Clipboard.Data)

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Enter PASSWORD")
	s.Assert().Contains(io.ReadLine(), "[=] Loaded Record: "+s.Record.ID.String())
	s.Assert().Contains(io.ReadLine(), "[=] PASSWORD copied to clipboard")
}

func (s *CopyTestSuite) TestRunPassUsernameCopied() {
	//-- arrange
	io := pipe.OpenStdio(1, 4, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", s.Password)
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name, "-username")

	//-- assert
	s.RequireResultPass()

	s.Assert().Equal(s.Record.Username, s.Clipboard.Data)

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Enter PASSWORD")
	s.Assert().Contains(io.ReadLine(), "[=] Loaded Record: "+s.Record.ID.String())
	s.Assert().Contains(io.ReadLine(), "[=] USERNAME copied to clipboard")
}

func (s *CopyTestSuite) TestRunQRFailsWithErrorRenderingToQRCode() {
	//-- arrange
	s.QRCode.RenderError = errors.New("")

	io := pipe.OpenStdio(1, 3, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", s.Password)
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name, "-qr")

	//-- assert
	s.RequireResultFail("error rendering qr-code")

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Enter PASSWORD")
	s.Assert().Contains(io.ReadLine(), "[=] Loaded Record: "+s.Record.ID.String())
}

func (s *CopyTestSuite) TestRunQRPassesAndRendersPasswordAsQRCode() {
	//-- arrange
	io := pipe.OpenStdio(1, 4, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", s.Password)
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name, "-qr")

	//-- assert
	s.RequireResultPass()

	s.Assert().Equal(s.Record.Password, s.QRCode.Text)

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Enter PASSWORD")
	s.Assert().Contains(io.ReadLine(), "[=] Loaded Record: "+s.Record.ID.String())
	s.Assert().Contains(io.ReadLine(), "[=] PASSWORD rendered as QR-Code")
}

func (s *CopyTestSuite) TestRunQRPassesAndRendersUsernameAsQRCode() {
	//-- arrange
	io := pipe.OpenStdio(1, 4, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", s.Password)
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name, "-username", "-qr")

	//-- assert
	s.RequireResultPass()

	s.Assert().Equal(s.Record.Username, s.QRCode.Text)

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Enter PASSWORD")
	s.Assert().Contains(io.ReadLine(), "[=] Loaded Record: "+s.Record.ID.String())
	s.Assert().Contains(io.ReadLine(), "[=] USERNAME rendered as QR-Code")
}
