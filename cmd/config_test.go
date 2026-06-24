package cmd_test

import (
	"pvault/cmd"
	"pvault/config"
	"pvault/errors"
	"testing"

	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	test.CommandSuite[*cmd.ConfigCommand]
	Loader *config.MockLoader[config.Config]
}

func TestConfigCommandSuite(t *testing.T) {
	s := ConfigTestSuite{
		Loader: &config.MockLoader[config.Config]{},
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewConfigCommand(s.Loader))
	suite.Run(t, &s)
}

func (s *ConfigTestSuite) SetupTest() {
	*s.Loader = config.MockLoader[config.Config]{}
}

func (s *ConfigTestSuite) TestRunFailErrorLoadingConfig() {
	//-- arrange
	s.Loader.Error = errors.New("")

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("error loading config")
}

func (s *ConfigTestSuite) TestRunValidConfigPrintsConfig() {
	//-- arrange
	r := rand.New(0)

	NAME := r.ASCII(15)
	CONFIG := config.Config{
		Version:    r.ASCII(3),
		VaultPath:  r.ASCII(15),
		OutputPath: r.ASCII(15),
	}

	s.Loader.Name = NAME
	s.Loader.Config = CONFIG

	out := pipe.OpenStdout(5)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()

	s.Assert().Contains(out.ReadLine(), "Loaded from "+NAME)
	s.Assert().Contains(out.ReadLine(), "Version: "+CONFIG.Version)
	out.SkipLines(1)
	s.Assert().Contains(out.ReadLine(), "Vault Path: "+CONFIG.VaultPath)
	s.Assert().Contains(out.ReadLine(), "Output Path: "+CONFIG.OutputPath)
}
