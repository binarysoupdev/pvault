package cmd_test

import (
	"path/filepath"
	"pvault/cmd"
	"pvault/config"
	"pvault/data"
	"pvault/vault"
	"pvault/vault/record"
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
	Record     record.Record
}

func TestLockCommandSuite(t *testing.T) {
	suite.Run(t, &LockTestSuite{
		CommandSuite: test.NewCommandSuite(cmd.NewLockCommand()),
	})
}

func (s *LockTestSuite) SetupTest() {
	config.SetGlobal(config.Config{
		VaultPath: file.NewPath(s.T(), "vault"),
	})

	err := vault.InitializeNew(config.Global.VaultPath)
	s.Require().NoError(err)

	rand := rand.New(0)
	s.RecordPath = file.NewPath(s.T(), rand.ASCII(10))
	s.Record = record.NewFromName(rand.ASCII(15))

	err = data.SaveJSON(s.Record, s.RecordPath)
	s.Require().NoError(err)
}

//=====================================

func (s *LockTestSuite) TestRunPathNotEmpty() {
	//-- act
	s.RunCommand("-path", "")

	//-- assert
	s.RequireResultFail("\"path\" cannot be empty")
}

func (s *LockTestSuite) TestRunInvalidVaultPath() {
	//-- arrange
	config.Global.VaultPath = "invalid"

	//-- act
	s.RunCommand("-path", s.RecordPath)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *LockTestSuite) TestRunInvalidRecordPath() {
	//-- act
	s.RunCommand("-path", "invalid")

	//-- assert
	s.RequireResultFail("error loading source record")
}

func (s *LockTestSuite) TestRunSaveNewValid() {
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
	s.Assert().Contains(out.ReadLine(), "[-] "+s.RecordPath)

	s.Assert().FileExists(VAULT_FILE)
	s.Assert().NoFileExists(s.RecordPath)
}

func (s *LockTestSuite) TestRunUpdateExistingValid() {
	//-- arrange
	VAULT_FILE := filepath.Join(config.Global.VaultPath, s.Record.ID.String()+".json")

	v, err := vault.Open(config.Global.VaultPath)
	s.Require().NoError(err)

	err = v.SaveRecord(record.Record{
		ID:   s.Record.ID,
		Name: "existing",
	})
	s.Require().NoError(err)

	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.RunCommand("-path", s.RecordPath)

	//-- assert
	s.RequireResultPass()

	line := out.ReadLine()
	s.Require().Contains(line, "[+] Updated Record: "+s.Record.ID.String())

	record, err := data.LoadJSON[record.Record](VAULT_FILE)
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, record)

	s.Assert().Contains(out.ReadLine(), "[-] "+s.RecordPath)
	s.Assert().NoFileExists(s.RecordPath)
}

func (s *LockTestSuite) TestRunExistingNameInvalid() {
	//-- arrange
	v, err := vault.Open(config.Global.VaultPath)
	s.Require().NoError(err)

	err = v.SaveRecord(record.NewFromName(s.Record.Name))
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-path", s.RecordPath)

	//-- assert
	s.RequireResultFail("error saving vault record")
}
