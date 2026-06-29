package cmd_test

import (
	"errors"
	"os"
	"pvault/cmd"
	"pvault/config"
	"pvault/data"
	"pvault/tools/clipboard"
	"pvault/vault"
	"pvault/vault/record"
	"testing"

	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type CopyTestSuite struct {
	test.CommandSuite[*cmd.CopyCommand]
	Clipboard *clipboard.MockClipboard

	ConfigLoader config.Loader[config.Config]
	Config       config.Config

	Vault    vault.Vault
	Record   record.Record
	Password string
}

func TestCopyCommandSuite(t *testing.T) {
	s := CopyTestSuite{
		ConfigLoader: config.NewLoader[config.Config](file.NewPath(t, "config.json")),
		Clipboard:    clipboard.Mock(),
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewCopyCommand(s.ConfigLoader, s.Clipboard))
	suite.Run(t, &s)
}

func (s *CopyTestSuite) SetupTest() {
	*s.Clipboard = clipboard.MockClipboard{}

	s.Config = config.Config{
		Version:   config.VERSION,
		VaultPath: file.NewPath(s.T(), "vault"),
	}
	err := data.SaveJSON(s.Config, s.ConfigLoader.ConfigPath)
	s.Require().NoError(err)

	rand := rand.New(0)
	s.Record = record.NewFromName(rand.ASCII(15))
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
	err := os.Remove(s.ConfigLoader.ConfigPath)
	s.Require().NoError(err)

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("error loading config")
}

func (s *CopyTestSuite) TestRunFailInvalidVaultPath() {
	//-- arrange
	s.Config.VaultPath = "invalid"
	err := data.SaveJSON(s.Config, s.ConfigLoader.ConfigPath)
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

func (s *CopyTestSuite) TestRunFailVaultFileMissing() {
	//-- arrange
	err := s.Vault.Database.DeleteRecord(s.Record.ID)
	s.Require().NoError(err)

	io := pipe.OpenStdio(1, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", s.Password)
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error reading record file")

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Enter PASSWORD")
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
	s.RequireResultFail("error decrypting ciphertext")

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
