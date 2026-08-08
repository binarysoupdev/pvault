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
	record_v1 "pvault/app/vault/record/record/legacy/v1"
	record_v2 "pvault/app/vault/record/record/v2"
	"testing"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UnlockTestSuite struct {
	test.CommandSuite[*cmds.UnlockCommand]
	ConfigLoader json.Loader[config.Config]
	Config       config.Config

	Vault    vault.Vault
	Record   record_v2.Record
	Password string
}

func TestUnlockCommandSuite(t *testing.T) {
	s := UnlockTestSuite{
		ConfigLoader: json.NewLoader[config.Config](filepath.Join(t.TempDir(), "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmds.NewUnlockCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *UnlockTestSuite) SetupTest() {
	s.Config = config.Config{
		Version:    config.VERSION,
		VaultPath:  filepath.Join(s.T().TempDir(), "vault"),
		OutputPath: s.T().TempDir(),
	}
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	s.Record = record_v2.NewEmptyRecord("name")
	s.Password = "Password123!"

	var err error
	s.Vault, err = vault.InitializeNew(s.Config.VaultPath, "")
	s.Require().NoError(err)

	s.Require().NoError(s.Vault.SaveRecord(s.Record, s.Password))
}

//=====================================

func (s *UnlockTestSuite) TestRunFailsWhenConfigNotFound() {
	//-- arrange
	s.Require().NoError(os.Remove(s.ConfigLoader.Path))

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("invalid config path")
}

func (s *UnlockTestSuite) TestRunFailsWhenConfigVersionUnsupported() {
	//-- arrange
	s.Config.Version = config.VERSION + 1
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail(fmt.Sprintf("unsupported config version \"%d\"", s.Config.Version))
}

func (s *UnlockTestSuite) TestRunFailsWhenVaultOutOfDate() {
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
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail(fmt.Sprintf("vault (@v%d) out-of-date", DATABASE.GetVersion()))
}

func (s *UnlockTestSuite) TestRunFailsWhenConfigOutputPathInvalid() {
	//-- arrange
	s.Config.OutputPath = "invalid"
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	//-- act
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error validating config \"output_path\"")
}

func (s *UnlockTestSuite) TestRunFailsWhenNoResults() {
	//-- act
	s.RunCommand("-s", "no match")

	//-- assert
	s.RequireResultFail("no matches found")
}

func (s *UnlockTestSuite) TestRunFailsWithIncorrectPassword() {
	//-- arrange
	io := pipe.OpenStdio(1, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", s.Password+"x")
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error loading vault record")

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Enter PASSWORD")
}

func (s *UnlockTestSuite) TestRunPassesAndUnlocksRecord() {
	//-- arrange
	OUTPUT_FILE := filepath.Join(s.Config.OutputPath, s.Record.ID.String()+".json")

	io := pipe.OpenStdio(1, 4, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", s.Password)
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultPass()

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Enter PASSWORD")
	s.Assert().Contains(io.ReadLine(), "[=] Loaded Record: "+s.Record.ID.String())
	s.Assert().Contains(io.ReadLine(), "[+] "+OUTPUT_FILE)

	record, err := json.UnmarshalFile[record_v2.Record](OUTPUT_FILE)
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, record)
}

func (s *UnlockTestSuite) TestRunPassesAndOutputRecordWasUpgradedWhenUnlockingOlderVersion() {
	//-- arrange
	r1 := record_v1.Record{
		ID:       uuid.New(),
		Name:     "record v1",
		Username: "foo",
		Password: "bar",
	}
	s.Require().NoError(s.Vault.SaveRecord(r1, s.Password))

	OUTPUT_FILE := filepath.Join(s.Config.OutputPath, r1.ID.String()+".json")

	io := pipe.OpenStdio(1, 4, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", s.Password)
	io.EndQueue()

	s.RunCommand("-s", r1.Name)

	//-- assert
	s.RequireResultPass()

	s.Assert().Contains(io.ReadLine(), r1.Name)
	s.Assert().Contains(io.ReadLine(), "Enter PASSWORD")
	s.Assert().Contains(io.ReadLine(), "[=] Loaded Record: "+r1.ID.String())
	s.Assert().Contains(io.ReadLine(), "[+] "+OUTPUT_FILE)

	record, err := json.UnmarshalFile[record_v2.Record](OUTPUT_FILE)
	s.Require().NoError(err)
	s.Assert().Equal(r1.Upgrade(), record)
}
