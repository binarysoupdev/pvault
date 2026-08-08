package cmds_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"pvault/app/cmds"
	"pvault/app/config"
	"pvault/tools/clipboard"
	"pvault/tools/qrcode"
	"pvault/vault"
	"pvault/vault/database"
	db_v1 "pvault/vault/database/encoder/legacy/v1"
	"pvault/vault/index"
	"pvault/vault/meta"
	meta_v1 "pvault/vault/meta/encoder/v1"
	record_v2 "pvault/vault/record/record/v2"
	"testing"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/suite"
)

type CopyTestSuite struct {
	test.CommandSuite[*cmds.CopyCommand]
	Clipboard *clipboard.MockClipboard
	QRCode    *qrcode.MockRenderer

	ConfigLoader json.Loader[config.Config]
	Config       config.Config

	Vault    vault.Vault
	Record   record_v2.Record
	Password string
}

func TestCopyCommandSuite(t *testing.T) {
	s := CopyTestSuite{
		ConfigLoader: json.NewLoader[config.Config](filepath.Join(t.TempDir(), "config.json")),
		Clipboard:    clipboard.Mock(),
		QRCode:       qrcode.Mock(),
	}

	s.CommandSuite = test.NewCommandSuite(cmds.NewCopyCommand(s.ConfigLoader, s.Clipboard, s.QRCode))
	suite.Run(t, &s)
}

func (s *CopyTestSuite) SetupTest() {
	*s.Clipboard = clipboard.MockClipboard{}
	*s.QRCode = qrcode.MockRenderer{}

	s.Config = config.Config{
		Version:   config.VERSION,
		VaultPath: filepath.Join(s.T().TempDir(), "vault"),
	}
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	s.Record = record_v2.NewEmptyRecord("name")
	s.Record.Username = "username"
	s.Record.Password = "password"

	s.Password = "Password123!"

	var err error
	s.Vault, err = vault.InitializeNew(s.Config.VaultPath, "")
	s.Require().NoError(err)

	s.Require().NoError(s.Vault.SaveRecord(s.Record, s.Password))
}

//=====================================

func (s *CopyTestSuite) TestRunFailsWhenClipboardUnsupported() {
	//-- arrange
	s.Clipboard.UnsupportedError = errors.New("")

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("clipboard unsupported")
}

func (s *CopyTestSuite) TestRunFailsWhenConfigNotFound() {
	//-- arrange
	s.Require().NoError(os.Remove(s.ConfigLoader.Path))

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("invalid config path")
}

func (s *CopyTestSuite) TestRunFailsWhenConfigVersionUnsupported() {
	//-- arrange
	s.Config.Version = config.VERSION + 1
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail(fmt.Sprintf("unsupported config version \"%d\"", s.Config.Version))
}

func (s *CopyTestSuite) TestRunFailsWithInvalidVaultPath() {
	//-- arrange
	s.Config.VaultPath = "invalid"
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	//-- act
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *CopyTestSuite) TestRunFailsWhenVaultOutOfDate() {
	//-- arrange
	s.Config.VaultPath = s.T().TempDir()
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	DATABASE := db_v1.Encoder{}

	META := meta.Metadata{
		DatabaseVersion: DATABASE.GetVersion(),
	}
	s.Require().NoError(meta.SaveMetadata(meta_v1.Encoder{}, s.Config.VaultPath, META))
	s.Require().NoError(database.SaveIndex(DATABASE, s.Config.VaultPath, index.IndexMap{}))

	//-- act
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail(fmt.Sprintf("vault (@v%d) out-of-date", DATABASE.GetVersion()))
}

func (s *CopyTestSuite) TestRunFailsWithNoResults() {
	//-- act
	s.RunCommand("-s", "no match")

	//-- assert
	s.RequireResultFail("no matches found")
}

func (s *CopyTestSuite) TestRunFailsWithIncorrectPassword() {
	//-- arrange
	io := pipe.OpenStdio(1, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", s.Password+"x")
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error loading vault record")

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Enter PASSWORD")
}

func (s *CopyTestSuite) TestRunFailsWhenErrorCopyingToClipboard() {
	//-- arrange
	s.Clipboard.WriteError = errors.New("")

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

func (s *CopyTestSuite) TestRunPassesAndCopiesPassword() {
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

func (s *CopyTestSuite) TestRunPassesAndCopiesUsername() {
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
