package cmd_test

import (
	"os"
	vault "pvault/app/vault/local"
	record "pvault/app/vault/record/version2"
	"pvault/cmd"
	"pvault/config"
	"testing"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type SearchTestSuite struct {
	test.CommandSuite[*cmd.SearchCommand]
	ConfigLoader json.Loader[config.Config]
	Config       config.Config
}

func TestSearchCommandSuite(t *testing.T) {
	s := SearchTestSuite{
		ConfigLoader: json.NewLoader[config.Config](file.NewPath(t, "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewSearchCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *SearchTestSuite) SetupTest() {
	s.Config = config.Config{
		Version:    config.VERSION,
		VaultPath:  file.NewPath(s.T(), "vault"),
		OutputPath: file.NewPath(s.T(), ""),
	}
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	_, err = vault.InitializeNew(s.Config.VaultPath)
	s.Require().NoError(err)
}

//=====================================

func (s *SearchTestSuite) TestRunFailErrorLoadingConfig() {
	//-- arrange
	err := os.Remove(s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("invalid config path")
}

func (s *SearchTestSuite) TestRunInvalidVaultPath() {
	//-- arrange
	rand := rand.New(0)
	NAME := rand.ASCII(15)

	s.Config.VaultPath = "invalid"
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-s", NAME)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *SearchTestSuite) TestRunInvalidNoResults() {
	//-- act
	s.RunCommand("-s", "no match")

	//-- assert
	s.RequireResultFail("no matches found")
}

func (s *SearchTestSuite) TestRunValidDisplaySearchResults() {
	//-- arrange
	rand := rand.New(0)
	NAME := rand.ASCII(15)

	v, err := vault.Open(s.Config.VaultPath)
	s.Require().NoError(err)

	err = v.SaveRecord(record.NewEmptyRecord(NAME), rand.ASCII(30))
	s.Require().NoError(err)

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-s", NAME)

	//-- assert
	s.RequireResultPass()

	line := out.ReadLine()
	s.Assert().Contains(line, "[1]")
	s.Assert().Contains(line, NAME)
}
