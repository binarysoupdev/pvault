package cmd_test

import (
	"os"
	"pvault/cmd"
	"pvault/config"
	"pvault/vault"
	"pvault/vault/record"
	"testing"

	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type DeleteTestSuite struct {
	test.CommandSuite[*cmd.DeleteCommand]
	Vault  vault.Vault
	Record record.Record
}

func TestDeleteCommandSuite(t *testing.T) {
	suite.Run(t, &DeleteTestSuite{
		CommandSuite: test.NewCommandSuite(cmd.NewDeleteCommand()),
	})
}

func (s *DeleteTestSuite) SetupTest() {
	config.SetGlobal(config.Config{
		VaultPath: file.NewPath(s.T(), "vault"),
	})

	rand := rand.New(0)
	s.Record = record.NewFromName(rand.ASCII(15))

	err := vault.InitializeNew(config.Global.VaultPath)
	s.Require().NoError(err)

	s.Vault, err = vault.Open(config.Global.VaultPath)
	s.Require().NoError(err)

	err = s.Vault.SaveRecord(s.Record)
	s.Require().NoError(err)
}

//=====================================

func (s *DeleteTestSuite) TestRunInvalidVaultPath() {
	//-- arrange
	config.Global.VaultPath = "invalid"

	//-- act
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *DeleteTestSuite) TestRunInvalidNoResults() {
	//-- act
	s.RunCommand("-s", "no match")

	//-- assert
	s.RequireResultFail("no matches found")
}

func (s *DeleteTestSuite) TestRunIncorrectConfirmName() {
	//-- arrange
	err := os.Remove(s.Vault.RecordPath(s.Record.ID))
	s.Require().NoError(err)

	io := pipe.OpenStdio(1, 2, true)
	defer io.Close()

	//-- act
	io.Queue("NAME: ", s.Record.Name+"x")
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("names do not match")

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Confirm NAME: "+s.Record.Name+"x")
}

func (s *DeleteTestSuite) TestRunVaultFileMissing() {
	//-- arrange
	err := os.Remove(s.Vault.RecordPath(s.Record.ID))
	s.Require().NoError(err)

	io := pipe.OpenStdio(1, 2, true)
	defer io.Close()

	//-- act
	io.Queue("NAME: ", s.Record.Name)
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error deleting vault record")

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Confirm NAME: "+s.Record.Name)
}

func (s *DeleteTestSuite) TestRunValid() {
	//-- arrange
	io := pipe.OpenStdio(1, 3, true)
	defer io.Close()

	//-- act
	io.Queue("NAME: ", s.Record.Name)
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultPass()
	s.Assert().NoFileExists(s.Vault.RecordPath(s.Record.ID))

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Confirm NAME: "+s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "[-] Deleted Record: "+s.Record.ID.String())
}
