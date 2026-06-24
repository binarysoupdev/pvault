package cmd_test

import (
	"path/filepath"
	"pvault/cmd"
	"pvault/config"
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

type CreateTestSuite struct {
	test.CommandSuite[*cmd.CreateCommand]
	ConfigLoader *config.MockLoader[config.Config]
}

func TestCreateCommandSuite(t *testing.T) {
	s := CreateTestSuite{
		ConfigLoader: &config.MockLoader[config.Config]{},
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewCreateCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *CreateTestSuite) SetupTest() {
	*s.ConfigLoader = config.MockLoader[config.Config]{
		Config: config.Config{
			VaultPath:  file.NewPath(s.T(), "vault"),
			OutputPath: file.NewPath(s.T(), ""),
		},
	}

	_, err := vault.InitializeNew(s.ConfigLoader.Config.VaultPath)
	s.Require().NoError(err)
}

//=====================================

func (s *CreateTestSuite) TestRunFailErrorLoadingConfig() {
	//-- arrange
	s.ConfigLoader.Error = errors.New("")

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("error loading config")
}

func (s *CreateTestSuite) TestRunFailConfigOutputPathInvalid() {
	//-- arrange
	s.ConfigLoader.Config.OutputPath = "invalid"

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("error validating config \"output_path\"")
}

func (s *CreateTestSuite) TestRunNameNotEmpty() {
	//-- act
	s.RunCommand("-name", "")

	//-- assert
	s.RequireResultFail("\"name\" cannot be empty")
}

func (s *CreateTestSuite) TestRunInvalidVaultPath() {
	//-- arrange
	rand := rand.New(0)
	NAME := rand.ASCII(15)

	s.ConfigLoader.Config.VaultPath = "invalid"

	//-- act
	s.RunCommand("-name", NAME)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *CreateTestSuite) TestRunInvalidNameAlreadyExists() {
	//-- arrange
	rand := rand.New(0)
	NAME := rand.ASCII(15)

	v, err := vault.Open(s.ConfigLoader.Config.VaultPath)
	s.Require().NoError(err)

	err = v.SaveRecord(record.NewFromName(NAME), rand.ASCII(30))
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-name", NAME)

	//-- assert
	s.RequireResultFail("error validating record")
}

func (s *CreateTestSuite) TestRunIncorrectVerifyPassword() {
	//-- arrange
	rand := rand.New(0)
	NAME := rand.ASCII(15)
	PASSWORD := rand.ASCII(30)

	io := pipe.OpenStdio(2, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.Queue("PASSWORD: ", PASSWORD+"x")
	io.EndQueue()

	s.RunCommand("-name", NAME)

	//-- assert
	s.RequireResultFail("passwords do not match")
	s.Assert().Contains(io.ReadLine(), "New PASSWORD")
	s.Assert().Contains(io.ReadLine(), "Verify PASSWORD")
}

func (s *CreateTestSuite) TestRunValid() {
	//-- arrange
	rand := rand.New(0)
	NAME := rand.ASCII(15)
	PASSWORD := rand.ASCII(30)

	io := pipe.OpenStdio(2, 4, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.Queue("PASSWORD: ", PASSWORD)
	io.EndQueue()

	s.RunCommand("-name", NAME)

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(io.ReadLine(), "New PASSWORD")
	s.Assert().Contains(io.ReadLine(), "Verify PASSWORD")

	line := io.ReadLine()
	s.Require().Contains(line, "[+] New Record: ")

	ID := line[len(line)-36:]
	VAULT_FILE := filepath.Join(s.ConfigLoader.Config.VaultPath, ID)
	OUTPUT_FILE := filepath.Join(s.ConfigLoader.Config.OutputPath, ID+".json")

	s.Assert().Contains(io.ReadLine(), "[+] "+OUTPUT_FILE)
	s.Assert().FileExists(VAULT_FILE)
	s.Assert().FileExists(OUTPUT_FILE)
}
