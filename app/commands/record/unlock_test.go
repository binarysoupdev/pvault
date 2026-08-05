package record_test

import (
	"os"
	"path/filepath"
	cmd "pvault/app/commands/record"
	"pvault/app/config"
	"pvault/app/vault"
	record_v1 "pvault/app/vault/record/record/legacy/v1"
	record_v2 "pvault/app/vault/record/record/v2"
	"testing"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UnlockTestSuite struct {
	test.CommandSuite[*cmd.UnlockCommand]
	ConfigLoader json.Loader[config.Config]
	Config       config.Config

	Vault    vault.Vault
	Record   record_v2.Record
	Password string
}

func TestUnlockCommandSuite(t *testing.T) {
	s := UnlockTestSuite{
		ConfigLoader: json.NewLoader[config.Config](file.NewPath(t, "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewUnlockCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *UnlockTestSuite) SetupTest() {
	s.Config = config.Config{
		Version:    config.VERSION,
		VaultPath:  file.NewPath(s.T(), "vault"),
		OutputPath: file.NewPath(s.T(), ""),
	}
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	rand := rand.New(0)
	s.Record = record_v2.NewEmptyRecord(rand.ASCII(15))
	s.Password = rand.ASCII(30)

	s.Vault, err = vault.InitializeNew(s.Config.VaultPath, "")
	s.Require().NoError(err)

	err = s.Vault.SaveRecord(s.Record, s.Password)
	s.Require().NoError(err)
}

//=====================================

func (s *UnlockTestSuite) TestRunFailsWhenErrorLoadingConfig() {
	//-- arrange
	err := os.Remove(s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("invalid config path")
}

func (s *UnlockTestSuite) TestRunFailsWhenInvalidVaultPath() {
	//-- arrange
	s.Config.VaultPath = "invalid"
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *UnlockTestSuite) TestRunFailsWhenConfigOutputPathInvalid() {
	//-- arrange
	s.Config.OutputPath = "invalid"
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

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
	s.RequireResultFail("error decrypting record")

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

	err := s.Vault.SaveRecord(r1, s.Password)
	s.Require().NoError(err)

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
