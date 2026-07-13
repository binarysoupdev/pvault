package cmd_test

import (
	"os"
	"path/filepath"
	"pvault/cmd"
	"pvault/config"
	"pvault/vault"
	"pvault/vault/record"
	"testing"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type UnlockTestSuite struct {
	test.CommandSuite[*cmd.UnlockCommand]
	ConfigLoader json.Loader[config.Config]
	Config       config.Config

	Vault    vault.Vault
	Record   record.RecordV2
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
	s.Record = record.NewFromName(rand.ASCII(15))
	s.Password = rand.ASCII(30)

	s.Vault, err = vault.InitializeNew(s.Config.VaultPath)
	s.Require().NoError(err)

	err = s.Vault.SaveRecord(s.Record, s.Password)
	s.Require().NoError(err)
}

//=====================================

func (s *UnlockTestSuite) TestRunFailErrorLoadingConfig() {
	//-- arrange
	err := os.Remove(s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("invalid config path")
}

func (s *UnlockTestSuite) TestRunInvalidVaultPath() {
	//-- arrange
	s.Config.VaultPath = "invalid"
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *UnlockTestSuite) TestRunFailConfigOutputPathInvalid() {
	//-- arrange
	s.Config.OutputPath = "invalid"
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error validating config \"output_path\"")
}

func (s *UnlockTestSuite) TestRunInvalidNoResults() {
	//-- act
	s.RunCommand("-s", "no match")

	//-- assert
	s.RequireResultFail("no matches found")
}

func (s *UnlockTestSuite) TestRunIncorrectPassword() {
	//-- arrange
	io := pipe.OpenStdio(1, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", s.Password+"x")
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error decrypting ciphertext")

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Enter PASSWORD")
}

func (s *UnlockTestSuite) TestRunValid() {
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

	record, err := json.UnmarshalFile[record.RecordV2](OUTPUT_FILE)
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, record)
}
