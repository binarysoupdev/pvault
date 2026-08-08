package cmds_test

import (
	"fmt"
	"os"
	"path/filepath"
	"pvault/app/cmds"
	"pvault/app/config"
	"pvault/app/vault"
	"pvault/app/vault/database"
	db_v1 "pvault/app/vault/database/encoder/legacy/v1"
	"pvault/app/vault/index"
	"pvault/app/vault/meta"
	meta_v1 "pvault/app/vault/meta/encoder/v1"
	record_v2 "pvault/app/vault/record/record/v2"
	"testing"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/suite"
)

type LockTestSuite struct {
	test.CommandSuite[*cmds.LockCommand]
	ConfigLoader json.Loader[config.Config]
	Config       config.Config

	RecordPath string
	Record     record_v2.Record
}

func TestLockCommandSuite(t *testing.T) {
	s := LockTestSuite{
		ConfigLoader: json.NewLoader[config.Config](filepath.Join(t.TempDir(), "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmds.NewLockCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *LockTestSuite) SetupTest() {
	s.Config = config.Config{
		Version:    config.VERSION,
		VaultPath:  filepath.Join(s.T().TempDir(), "vault"),
		OutputPath: s.T().TempDir(),
	}
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	_, err := vault.InitializeNew(s.Config.VaultPath, "")
	s.Require().NoError(err)

	s.RecordPath = filepath.Join(s.T().TempDir(), "record.json")
	s.Record = record_v2.NewEmptyRecord("name")

	s.Require().NoError(json.MarshalFile(s.Record, s.RecordPath))
}

//=====================================

func (s *LockTestSuite) TestRunFailsWhenConfigNotFound() {
	//-- arrange
	s.Require().NoError(os.Remove(s.ConfigLoader.Path))

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("invalid config path")
}

func (s *LockTestSuite) TestRunFailsWhenConfigVersionUnsupported() {
	//-- arrange
	s.Config.Version = config.VERSION + 1
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail(fmt.Sprintf("unsupported config version \"%d\"", s.Config.Version))
}

func (s *LockTestSuite) TestRunFailsWithInvalidVaultPath() {
	//-- arrange
	s.Config.VaultPath = "invalid"
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	//-- act
	s.RunCommand("-path", s.RecordPath)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *LockTestSuite) TestRunFailsWhenVaultOutOfDate() {
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
	s.RunCommand("-path", s.RecordPath)

	//-- assert
	s.RequireResultFail(fmt.Sprintf("vault (@v%d) out-of-date", DATABASE.GetVersion()))
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

	s.Require().NoError(json.MarshalFile(s.Record, s.RecordPath))

	//-- act
	s.RunCommand("-path", s.RecordPath)

	//-- assert
	s.RequireResultFail("error validating record")
}

func (s *LockTestSuite) TestRunFailsWhenNameAlreadyExists() {
	//-- arrange
	const PASSWORD = "Password123!"

	v, err := vault.Open(s.Config.VaultPath)
	s.Require().NoError(err)

	s.Require().NoError(v.SaveRecord(record_v2.NewEmptyRecord(s.Record.Name), PASSWORD))

	//-- act
	s.RunCommand("-path", s.RecordPath)

	//-- assert
	s.RequireResultFail("error validating record")
}

func (s *LockTestSuite) TestRunFailsWithIncorrectVerifyPassword() {
	//-- arrange
	const PASSWORD = "Password123!"

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
	const PASSWORD = "Password123!"

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

	const PASSWORD = "Password123!"

	v, err := vault.Open(s.Config.VaultPath)
	s.Require().NoError(err)

	s.Require().NoError(v.SaveRecord(OLD_RECORD, PASSWORD+"x"))

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
