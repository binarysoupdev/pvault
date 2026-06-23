package cmd_test

import (
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

	_, err := vault.InitializeNew(config.Global.VaultPath)
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

func (s *LockTestSuite) TestRunIncorrectVerifyPassword() {
	//-- arrange
	rand := rand.New(0)
	PASSWORD := rand.ASCII(30)

	io := pipe.OpenStdio(2, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.Queue("PASSWORD: ", PASSWORD+"x")
	io.EndQueue()

	s.RunCommand("-path", s.RecordPath)

	//-- assert
	s.RequireResultFail("passwords do not match")
	s.Assert().Contains(io.ReadLine(), "New PASSWORD")
	s.Assert().Contains(io.ReadLine(), "Verify PASSWORD")
}

func (s *LockTestSuite) TestRunValidSaveNew() {
	//-- arrange
	rand := rand.New(0)
	PASSWORD := rand.ASCII(30)

	io := pipe.OpenStdio(2, 4, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.Queue("PASSWORD: ", PASSWORD)
	io.EndQueue()

	s.RunCommand("-path", s.RecordPath)

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(io.ReadLine(), "New PASSWORD")
	s.Assert().Contains(io.ReadLine(), "Verify PASSWORD")

	s.Require().Contains(io.ReadLine(), "[+] Updated Record: "+s.Record.ID.String())
	s.Assert().Contains(io.ReadLine(), "[-] "+s.RecordPath)
	s.Assert().NoFileExists(s.RecordPath)

	v, err := vault.Open(config.Global.VaultPath)
	s.Require().NoError(err)

	res, err := v.LoadRecord(s.Record.Name, PASSWORD)
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, res)
}

func (s *LockTestSuite) TestRunValidUpdateExisting() {
	//-- arrange
	OLD_RECORD := record.NewFromName(s.Record.Name + "x")
	OLD_RECORD.ID = s.Record.ID

	rand := rand.New(0)
	PASSWORD := rand.ASCII(30)

	v, err := vault.Open(config.Global.VaultPath)
	s.Require().NoError(err)

	err = v.SaveRecord(OLD_RECORD, PASSWORD+"x")
	s.Require().NoError(err)

	io := pipe.OpenStdio(2, 4, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.Queue("PASSWORD: ", PASSWORD)
	io.EndQueue()

	s.RunCommand("-path", s.RecordPath)

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(io.ReadLine(), "New PASSWORD")
	s.Assert().Contains(io.ReadLine(), "Verify PASSWORD")

	s.Assert().Contains(io.ReadLine(), "[+] Updated Record: "+s.Record.ID.String())
	s.Assert().Contains(io.ReadLine(), "[-] "+s.RecordPath)
	s.Assert().NoFileExists(s.RecordPath)

	v, err = vault.Open(v.Path)
	s.Require().NoError(err)

	res, err := v.LoadRecord(s.Record.Name, PASSWORD)
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, res)
}
