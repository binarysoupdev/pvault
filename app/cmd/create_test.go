package cmd_test

import (
	"os"
	"path/filepath"
	"pvault/app/cmd"
	"pvault/app/config"
	"pvault/app/vault/local"
	v2 "pvault/app/vault/record/version2"
	"testing"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type CreateTestSuite struct {
	test.CommandSuite[*cmd.CreateCommand]
	ConfigLoader json.Loader[config.Config]
	Config       config.Config

	Vault local.Vault
}

func TestCreateCommandSuite(t *testing.T) {
	s := CreateTestSuite{
		ConfigLoader: json.NewLoader[config.Config](file.NewPath(t, "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewCreateCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *CreateTestSuite) SetupTest() {
	s.Config = config.Config{
		Version:    config.VERSION,
		VaultPath:  file.NewPath(s.T(), "vault"),
		OutputPath: file.NewPath(s.T(), ""),
	}
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	s.Vault, err = local.CreateNewVault(s.Config.VaultPath)
	s.Require().NoError(err)
}

//=====================================

func (s *CreateTestSuite) TestRunFailErrorLoadingConfig() {
	//-- arrange
	err := os.Remove(s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("invalid config path")
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

	s.Config.VaultPath = "invalid"
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-name", NAME)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *CreateTestSuite) TestRunFailConfigOutputPathInvalid() {
	//-- arrange
	const NAME = "name"

	s.Config.OutputPath = "invalid"
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-name", NAME)

	//-- assert
	s.RequireResultFail("error validating config \"output_path\"")
}

func (s *CreateTestSuite) TestRunInvalidNameAlreadyExists() {
	//-- arrange
	rand := rand.New(0)
	NAME := rand.ASCII(15)

	v, err := local.OpenVault(s.Config.VaultPath)
	s.Require().NoError(err)

	err = v.SaveRecord(v2.NewEmptyRecord(NAME), rand.ASCII(30))
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
	s.Require().Contains(line, "[+] Saved Record: ")

	ID, err := uuid.Parse(line[len(line)-36:])
	s.Require().NoError(err)

	OUTPUT_FILE := filepath.Join(s.Config.OutputPath, ID.String()+".json")
	s.Assert().Contains(io.ReadLine(), "[+] "+OUTPUT_FILE)

	s.Vault.Map, err = s.Vault.Index.LoadMap()
	s.Require().NoError(err)

	r1, err := json.UnmarshalFile[v2.Record](OUTPUT_FILE)
	s.Require().NoError(err)

	r2, err := s.Vault.LoadRecord(NAME, PASSWORD)
	s.Require().NoError(err)

	s.Assert().Equal(r1, r2)
}
