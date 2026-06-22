package cmd_test

import (
	"os"
	"path/filepath"
	"pvault/cmd"
	"pvault/config"
	"pvault/data"
	"pvault/vault"
	"testing"

	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type UnlockTestSuite struct {
	test.CommandSuite[*cmd.UnlockCommand]
	Vault  vault.Vault
	Record vault.Record
}

func TestUnlockCommandSuite(t *testing.T) {
	suite.Run(t, &UnlockTestSuite{
		CommandSuite: test.NewCommandSuite(cmd.NewUnlockCommand()),
	})
}

func (s *UnlockTestSuite) SetupTest() {
	config.SetGlobal(config.Config{
		VaultPath:  file.NewPath(s.T(), "vault"),
		OutputPath: file.NewPath(s.T(), ""),
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

func (s *UnlockTestSuite) TestRunInvalidVaultPath() {
	//-- arrange
	config.Global.VaultPath = "invalid"

	//-- act
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *UnlockTestSuite) TestRunValidNoResults() {
	//-- arrange
	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-s", "no match")

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), "No MATCHES found")
}

func (s *UnlockTestSuite) TestRunVaultFileMissing() {
	//-- arrange
	err := os.Remove(s.Vault.RecordPath(s.Record.ID))
	s.Require().NoError(err)

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error loading vault record")
	s.Assert().Contains(out.ReadLine(), s.Record.Name)
}

func (s *UnlockTestSuite) TestRunInvalidOutputPath() {
	//-- arrange
	config.Global.OutputPath = "invalid"

	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error creating output record")
	s.Assert().Contains(out.ReadLine(), s.Record.Name)
	s.Assert().Contains(out.ReadLine(), "[=] Loaded Record: "+s.Record.ID.String())
}

func (s *UnlockTestSuite) TestRunValid() {
	//-- arrange
	OUTPUT_FILE := filepath.Join(config.Global.OutputPath, s.Record.ID.String()+".json")

	out := pipe.OpenStdout(3)
	defer out.Close()

	//-- act
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), s.Record.Name)

	line := out.ReadLine()
	s.Require().Contains(line, "[=] Loaded Record: "+s.Record.ID.String())

	record, err := data.LoadJSON[vault.Record](OUTPUT_FILE)
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, record)

	s.Assert().Contains(out.ReadLine(), "[+] "+OUTPUT_FILE)
}
