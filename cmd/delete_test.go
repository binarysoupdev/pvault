package cmd_test

import (
	"os"
	"pvault/cmd"
	"pvault/config"
	"pvault/vault"
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
	Record vault.Record
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
	s.Record = vault.EmptyRecord(rand.ASCII(15))

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

func (s *DeleteTestSuite) TestRunValidNoResults() {
	//-- arrange
	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-s", "no match")

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), "No MATCHES found")
}

func (s *DeleteTestSuite) TestRunVaultFileMissing() {
	//-- arrange
	err := os.Remove(s.Vault.RecordPath(s.Record.ID))
	s.Require().NoError(err)

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error deleting vault record")
	s.Assert().Contains(out.ReadLine(), s.Record.Name)
}

func (s *DeleteTestSuite) TestRunValid() {
	//-- arrange
	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultPass()
	s.Assert().NoFileExists(s.Vault.RecordPath(s.Record.ID))

	s.Assert().Contains(out.ReadLine(), s.Record.Name)
	s.Assert().Contains(out.ReadLine(), "[-] Deleted Record: "+s.Record.ID.String())
}
