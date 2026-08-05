package vault_test

import (
	"os"
	"path/filepath"
	cmd "pvault/app/commands/vault"
	"pvault/app/config"
	"pvault/app/vault"
	record_v2 "pvault/app/vault/record/record/v2"
	"testing"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/suite"
)

type SearchTestSuite struct {
	test.CommandSuite[*cmd.SearchCommand]
	ConfigLoader json.Loader[config.Config]
	Config       config.Config
}

func TestSearchCommandSuite(t *testing.T) {
	s := SearchTestSuite{
		ConfigLoader: json.NewLoader[config.Config](filepath.Join(t.TempDir(), "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewSearchCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *SearchTestSuite) SetupTest() {
	s.Config = config.Config{
		Version:    config.VERSION,
		VaultPath:  filepath.Join(s.T().TempDir(), "vault"),
		OutputPath: s.T().TempDir(),
	}
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	_, err = vault.InitializeNew(s.Config.VaultPath, "")
	s.Require().NoError(err)
}

//=====================================

func (s *SearchTestSuite) TestRunFailsWhenErrorLoadingConfig() {
	//-- arrange
	err := os.Remove(s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("invalid config path")
}

func (s *SearchTestSuite) TestRunFailsWithInvalidVaultPath() {
	//-- arrange
	const NAME = "name"

	s.Config.VaultPath = "invalid"
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-s", NAME)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *SearchTestSuite) TestRunFailsWhenNoResults() {
	//-- act
	s.RunCommand("-s", "no match")

	//-- assert
	s.RequireResultFail("no matches found")
}

func (s *SearchTestSuite) TestRunPassesAndDisplaysSearchResults() {
	//-- arrange
	const NAME = "name"
	const PASSWORD = "Password123!"

	v, err := vault.Open(s.Config.VaultPath)
	s.Require().NoError(err)

	err = v.SaveRecord(record_v2.NewEmptyRecord(NAME), PASSWORD)
	s.Require().NoError(err)

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-s", NAME)

	//-- assert
	s.RequireResultPass()

	line := out.ReadLine()
	s.Assert().Contains(line, "[0]")
	s.Assert().Contains(line, NAME)
}
