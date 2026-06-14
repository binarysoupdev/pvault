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
	"github.com/stretchr/testify/suite"
)

type LockTestSuite struct {
	test.CommandSuite[*cmd.LockCommand]
	RecordPath string
	Record     vault.Record
}

func TestLockCommandSuite(t *testing.T) {
	suite.Run(t, &LockTestSuite{
		CommandSuite: test.NewCommandSuite(cmd.NewLockCommand()),
	})
}

func (s *LockTestSuite) SetupTest() {
	config.SetGlobal(config.Config{
		VaultPath: file.NewPath(s.T(), ""),
	})

	rand := rand.New(0)
	s.RecordPath = file.NewPath(s.T(), rand.ASCII(10))
	s.Record = vault.EmptyRecord(rand.ASCII(15))

	err := data.SaveJSON(s.Record, s.RecordPath)
	s.Require().NoError(err)
}

//=====================================

func (s *LockTestSuite) TestPathNotEmpty() {
	//-- act
	s.RunCommand("-path", "")

	//-- assert
	s.RequireResultFail("\"path\" cannot be empty")
}

func (s *LockTestSuite) TestInvalidRecordPath() {
	//-- act
	s.RunCommand("-path", "invalid")

	//-- assert
	s.RequireResultFail("error loading source record")
}

func (s *LockTestSuite) TestInvalidVaultPath() {
	//-- arrange
	config.Global.VaultPath += "/invalid"

	//-- act
	s.RunCommand("-path", s.RecordPath)

	//-- assert
	s.RequireResultFail("error saving vault record")
}

func (s *LockTestSuite) TestLockRecord() {
	//-- arrange
	VAULT_FILE := filepath.Join(config.Global.VaultPath, s.Record.ID.String()+".json")

	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.RunCommand("-path", s.RecordPath)

	//-- assert
	s.RequireResultPass()

	line := out.ReadLine()
	s.Require().Contains(line, "[+] Updated Record: "+s.Record.ID.String())

	record, err := data.LoadJSON[vault.Record](VAULT_FILE)
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, record)

	s.Assert().Contains(out.ReadLine(), "[-] "+s.RecordPath)
	s.Assert().NoFileExists(s.RecordPath)
}
