package cmd_test

import (
	"pvault/cmd"
	"pvault/config"
	"pvault/errors"
	"testing"

	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/suite"
)

type InitTestSuite struct {
	test.CommandSuite[*cmd.InitCommand]
	ConfigLoader *config.MockLoader[config.Config]
}

func TestInitCommandSuite(t *testing.T) {
	s := InitTestSuite{
		ConfigLoader: &config.MockLoader[config.Config]{},
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewInitCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *InitTestSuite) SetupTest() {
	*s.ConfigLoader = config.MockLoader[config.Config]{}
}

//=====================================

func (s *InitTestSuite) TestRunFailErrorLoadingConfig() {
	//-- arrange
	s.ConfigLoader.Error = errors.New("")

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("error loading config")
}

func (s *InitTestSuite) TestRunInvalidVaultPath() {
	//-- arrange
	*s.ConfigLoader = config.MockLoader[config.Config]{
		Config: config.Config{
			VaultPath: file.NewPath(s.T(), ""),
		},
	}

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("error initializing new vault")
}

func (s *InitTestSuite) TestRunValid() {
	//-- arrange
	*s.ConfigLoader = config.MockLoader[config.Config]{
		Config: config.Config{
			VaultPath: file.NewPath(s.T(), "vault"),
		},
	}

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), "[+] New Vault Initialized: "+s.ConfigLoader.Config.VaultPath)
}
