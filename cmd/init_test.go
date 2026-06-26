package cmd_test

import (
	"pvault/cmd"
	"pvault/config"
	"pvault/data"
	"testing"

	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/suite"
)

type InitTestSuite struct {
	test.CommandSuite[*cmd.InitCommand]
	ConfigLoader config.Loader[config.Config]
}

func TestInitCommandSuite(t *testing.T) {
	s := InitTestSuite{
		ConfigLoader: config.NewLoader[config.Config](file.NewPath(t, "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewInitCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

//=====================================

func (s *InitTestSuite) TestRunFailConfigNotFound() {
	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("error loading config")
}

func (s *InitTestSuite) TestRunInvalidVaultPath() {
	//-- arrange
	CONFIG := config.Config{
		Version:   config.VERSION,
		VaultPath: file.NewPath(s.T(), ""),
	}
	err := data.SaveJSON(CONFIG, s.ConfigLoader.ConfigPath)
	s.Require().NoError(err)

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("error initializing new vault")
}

func (s *InitTestSuite) TestRunValid() {
	//-- arrange
	CONFIG := config.Config{
		Version:   config.VERSION,
		VaultPath: file.NewPath(s.T(), "vault"),
	}
	err := data.SaveJSON(CONFIG, s.ConfigLoader.ConfigPath)
	s.Require().NoError(err)

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), "[+] New Vault Initialized: "+s.ConfigLoader.Config.VaultPath)
}
