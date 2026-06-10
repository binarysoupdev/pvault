package cmd_test

import (
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
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UnlockTestSuite struct {
	test.CommandSuite[*cmd.UnlockCommand]
	Record vault.Record
}

func TestUnlockCommandSuite(t *testing.T) {
	suite.Run(t, &UnlockTestSuite{
		CommandSuite: test.NewCommandSuite(cmd.NewUnlockCommand()),
	})
}

func (s *UnlockTestSuite) SetupTest() {
	config.SetGlobal(config.Config{
		VaultPath:  file.NewPath(s.T(), ""),
		OutputPath: file.NewPath(s.T(), ""),
	})

	rand := rand.New(0)
	s.Record = vault.NewRecord(rand.ASCII(15))

	err := vault.Vault{}.SaveRecord(s.Record)
	s.Require().NoError(err)
}

//=====================================

func (s *UnlockTestSuite) TestNameNotEmpty() {
	//-- act
	s.RunCommand("-name", "")

	//-- assert
	s.RequireResultFail("\"name\" cannot be empty")
}

func (s *UnlockTestSuite) TestInvalidIDFormat() {
	//-- act
	s.RunCommand("-name", "invalid")

	//-- assert
	s.RequireResultFail("error parsing ID")
}

func (s *UnlockTestSuite) TestInvalidID() {
	//-- act
	s.RunCommand("-name", uuid.Nil.String())

	//-- assert
	s.RequireResultFail("error loading vault record")
}

func (s *UnlockTestSuite) TestInvalidVaultPath() {
	//-- arrange
	config.Global.VaultPath += "/invalid"

	//-- act
	s.RunCommand("-name", s.Record.ID.String())

	//-- assert
	s.RequireResultFail("error loading vault record")
}

func (s *UnlockTestSuite) TestInvalidOutputPath() {
	//-- arrange
	config.Global.OutputPath += "/invalid"

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-name", s.Record.ID.String())

	//-- assert
	s.RequireResultFail("error creating output record")
	s.Assert().Contains(out.ReadLine(), "[=] Loaded Record: "+s.Record.ID.String())
}

func (s *UnlockTestSuite) TestUnlockRecord() {
	//-- arrange
	OUTPUT_FILE := filepath.Join(config.Global.OutputPath, s.Record.ID.String()+".json")

	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.RunCommand("-name", s.Record.ID.String())

	//-- assert
	s.RequireResultPass()

	line := out.ReadLine()
	s.Require().Contains(line, "[=] Loaded Record: "+s.Record.ID.String())

	record, err := data.LoadJSON[vault.Record](OUTPUT_FILE)
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, record)

	s.Assert().Contains(out.ReadLine(), "[+] "+OUTPUT_FILE)
}
