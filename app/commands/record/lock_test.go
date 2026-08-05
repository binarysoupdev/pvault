package record_test

import (
	"os"
	cmd "pvault/app/commands/record"
	"pvault/app/config"
	"pvault/app/vault"
	record_v2 "pvault/app/vault/record/record/v2"
	"testing"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type LockTestSuite struct {
	test.CommandSuite[*cmd.LockCommand]
	ConfigLoader json.Loader[config.Config]
	Config       config.Config

	RecordPath string
	Record     record_v2.Record
}

func TestLockCommandSuite(t *testing.T) {
	s := LockTestSuite{
		ConfigLoader: json.NewLoader[config.Config](file.NewPath(t, "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewLockCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *LockTestSuite) SetupTest() {
	s.Config = config.Config{
		Version:    config.VERSION,
		VaultPath:  file.NewPath(s.T(), "vault"),
		OutputPath: file.NewPath(s.T(), ""),
	}
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	_, err = vault.InitializeNew(s.Config.VaultPath, "")
	s.Require().NoError(err)

	rand := rand.New(0)
	s.RecordPath = file.NewPath(s.T(), rand.ASCII(10))
	s.Record = record_v2.NewEmptyRecord(rand.ASCII(15))

	err = json.MarshalFile(s.Record, s.RecordPath)
	s.Require().NoError(err)
}

//=====================================

func (s *LockTestSuite) TestRunFailsWhenErrorLoadingConfig() {
	//-- arrange
	err := os.Remove(s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("invalid config path")
}

func (s *LockTestSuite) TestRunFailsWhenPathIsEmpty() {
	//-- act
	s.RunCommand("-path", "")

	//-- assert
	s.RequireResultFail("\"path\" cannot be empty")
}

func (s *LockTestSuite) TestRunFailsWithInvalidVaultPath() {
	//-- arrange
	s.Config.VaultPath = "invalid"
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-path", s.RecordPath)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *LockTestSuite) TestRunFailsWithInvalidRecordPath() {
	//-- act
	s.RunCommand("-path", "invalid")

	//-- assert
	s.RequireResultFail("error loading source record")
}

func (s *LockTestSuite) TestRunFailsWithInvalidRecord() {
	//-- arrange
	s.Record.Name = ""

	err := json.MarshalFile(s.Record, s.RecordPath)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-path", s.RecordPath)

	//-- assert
	s.RequireResultFail("error validating record")
}

func (s *LockTestSuite) TestRunFailsWhenNameAlreadyExists() {
	//-- arrange
	v, err := vault.Open(s.Config.VaultPath)
	s.Require().NoError(err)

	rand := rand.New(0)
	err = v.SaveRecord(record_v2.NewEmptyRecord(s.Record.Name), rand.ASCII(30))
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-path", s.RecordPath)

	//-- assert
	s.RequireResultFail("error validating record")
}

func (s *LockTestSuite) TestRunFailsWithIncorrectVerifyPassword() {
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

func (s *LockTestSuite) TestRunPassesAndSavesNewRecord() {
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

	s.Require().Contains(io.ReadLine(), "[+] Saved Record: "+s.Record.ID.String())
	s.Assert().Contains(io.ReadLine(), "[-] "+s.RecordPath)
	s.Assert().NoFileExists(s.RecordPath)

	v, err := vault.Open(s.Config.VaultPath)
	s.Require().NoError(err)

	res, err := v.LoadRecord(s.Record.Name, PASSWORD)
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, res)
}

func (s *LockTestSuite) TestRunPassesAndUpdatesExistingRecord() {
	//-- arrange
	OLD_RECORD := record_v2.NewEmptyRecord(s.Record.Name + "x")
	OLD_RECORD.ID = s.Record.ID

	rand := rand.New(0)
	PASSWORD := rand.ASCII(30)

	v, err := vault.Open(s.Config.VaultPath)
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

	s.Assert().Contains(io.ReadLine(), "[+] Saved Record: "+s.Record.ID.String())
	s.Assert().Contains(io.ReadLine(), "[-] "+s.RecordPath)
	s.Assert().NoFileExists(s.RecordPath)

	v, err = vault.Open(v.Path)
	s.Require().NoError(err)

	res, err := v.LoadRecord(s.Record.Name, PASSWORD)
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, res)
}
