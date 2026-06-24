package cmd_test

import (
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

type SearchTestSuite struct {
	test.CommandSuite[*cmd.SearchCommand]
	ConfigLoader *config.MockLoader[config.Config]
}

func TestSearchCommandSuite(t *testing.T) {
	s := SearchTestSuite{
		ConfigLoader: &config.MockLoader[config.Config]{},
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewSearchCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *SearchTestSuite) SetupTest() {
	*s.ConfigLoader = config.MockLoader[config.Config]{
		Config: config.Config{
			VaultPath: file.NewPath(s.T(), "vault"),
		},
	}

	_, err := vault.InitializeNew(s.ConfigLoader.Config.VaultPath)
	s.Require().NoError(err)
}

//=====================================

func (s *SearchTestSuite) TestRunFailErrorLoadingConfig() {
	//-- arrange
	s.ConfigLoader.Error = errors.New("")

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("error loading config")
}

func (s *SearchTestSuite) TestRunInvalidVaultPath() {
	//-- arrange
	rand := rand.New(0)
	NAME := rand.ASCII(15)

	s.ConfigLoader.Config.VaultPath = "invalid"

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

	v, err := vault.Open(s.ConfigLoader.Config.VaultPath)
	s.Require().NoError(err)

	err = v.SaveRecord(record.NewFromName(NAME), rand.ASCII(30))
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
