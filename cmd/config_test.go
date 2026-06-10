package cmd_test

import (
	"pvault/cmd"
	"pvault/config"
	"testing"

	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	test.CommandSuite[*cmd.ConfigCommand]
}

func TestConfigCommandSuite(t *testing.T) {
	suite.Run(t, &ConfigTestSuite{
		CommandSuite: test.NewCommandSuite(cmd.NewConfigCommand()),
	})
}

func (s *ConfigTestSuite) TestPrintConfig() {
	//-- arrange
	r := rand.New(0)

	cfg := config.Config{
		Path:       r.ASCII(15),
		Version:    r.ASCII(3),
		VaultPath:  r.ASCII(15),
		OutputPath: r.ASCII(15),
	}
	config.SetGlobal(cfg)

	out := pipe.OpenStdout(5)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()

	s.Assert().Contains(out.ReadLine(), "Loaded from "+cfg.Path)
	s.Assert().Contains(out.ReadLine(), "Version: "+cfg.Version)
	out.SkipLines(1)
	s.Assert().Contains(out.ReadLine(), "Vault Path: "+cfg.VaultPath)
	s.Assert().Contains(out.ReadLine(), "Output Path: "+cfg.OutputPath)
}
