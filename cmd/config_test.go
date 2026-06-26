package cmd_test

import (
	"fmt"
	"pvault/cmd"
	"pvault/config"
	"pvault/data"
	"testing"

	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	test.CommandSuite[*cmd.ConfigCommand]
	ConfigLoader config.Loader[config.Config]
}

func TestConfigCommandSuite(t *testing.T) {
	s := ConfigTestSuite{
		ConfigLoader: config.NewLoader[config.Config](file.NewPath(t, "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewConfigCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

//=====================================

func (s *ConfigTestSuite) TestRunFailConfigNotFound() {
	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("error loading config")
}

func (s *ConfigTestSuite) TestRunValidConfigPrintsConfig() {
	//-- arrange
	rand := rand.New(0)

	CONFIG := config.Config{
		Version:    config.VERSION,
		VaultPath:  rand.ASCII(15),
		OutputPath: rand.ASCII(15),
	}
	err := data.SaveJSON(CONFIG, s.ConfigLoader.ConfigPath)
	s.Require().NoError(err)

	out := pipe.OpenStdout(5)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()

	s.Assert().Contains(out.ReadLine(), "Loaded from "+s.ConfigLoader.ConfigPath)
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Version [%d]", CONFIG.Version))
	out.SkipLines(1)
	s.Assert().Contains(out.ReadLine(), "Vault Path: "+CONFIG.VaultPath)
	s.Assert().Contains(out.ReadLine(), "Output Path: "+CONFIG.OutputPath)
}
