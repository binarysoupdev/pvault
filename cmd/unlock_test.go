package cmd_test

import (
	"os"
	"path/filepath"
	"pvault/cmd"
	"pvault/config"
	"pvault/data"
	"pvault/errors"
	"pvault/vault"
	"pvault/vault/record"
	"testing"

	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type UnlockTestSuite struct {
	test.CommandSuite[*cmd.UnlockCommand]
	ConfigLoader *config.MockLoader[config.Config]

	Vault    vault.Vault
	Record   record.Record
	Password string
}

func TestUnlockCommandSuite(t *testing.T) {
	s := UnlockTestSuite{
		ConfigLoader: &config.MockLoader[config.Config]{},
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewUnlockCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *UnlockTestSuite) SetupTest() {
	*s.ConfigLoader = config.MockLoader[config.Config]{
		Config: config.Config{
			VaultPath:  file.NewPath(s.T(), "vault"),
			OutputPath: file.NewPath(s.T(), ""),
		},
	}

	rand := rand.New(0)
	s.Record = record.NewFromName(rand.ASCII(15))
	s.Password = rand.ASCII(30)

	var err error
	s.Vault, err = vault.InitializeNew(s.ConfigLoader.Config.VaultPath)
	s.Require().NoError(err)

	err = s.Vault.SaveRecord(s.Record, s.Password)
	s.Require().NoError(err)
}

//=====================================

func (s *UnlockTestSuite) TestRunFailErrorLoadingConfig() {
	//-- arrange
	s.ConfigLoader.Error = errors.New("")

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("error loading config")
}

func (s *UnlockTestSuite) TestRunFailConfigOutputPathInvalid() {
	//-- arrange
	s.ConfigLoader.Config.OutputPath = "invalid"

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("error validating config \"output_path\"")
}

func (s *UnlockTestSuite) TestRunInvalidVaultPath() {
	//-- arrange
	s.ConfigLoader.Config.VaultPath = "invalid"

	//-- act
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *UnlockTestSuite) TestRunInvalidNoResults() {
	//-- act
	s.RunCommand("-s", "no match")

	//-- assert
	s.RequireResultFail("no matches found")
}

func (s *UnlockTestSuite) TestRunVaultFileMissing() {
	//-- arrange
	err := os.Remove(s.Vault.RecordPath(s.Record.ID))
	s.Require().NoError(err)

	io := pipe.OpenStdio(1, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", s.Password)
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error reading record file")

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Enter PASSWORD")
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
	OUTPUT_FILE := filepath.Join(s.ConfigLoader.Config.OutputPath, s.Record.ID.String()+".json")

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

	record, err := data.LoadJSON[record.Record](OUTPUT_FILE)
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, record)
}
