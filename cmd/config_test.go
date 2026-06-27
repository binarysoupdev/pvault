package cmd_test

import (
	"fmt"
	"os"
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

func (s *ConfigTestSuite) TestRunValidateConfigPrintsConfig() {
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

	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Loaded from \"%s\"", s.ConfigLoader.ConfigPath))
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Version [%d]", CONFIG.Version))
	out.SkipLines(1)

	//TODO: test prints for valid/invalid values
	s.Assert().Contains(out.ReadLine(), CONFIG.VaultPath)
	s.Assert().Contains(out.ReadLine(), CONFIG.OutputPath)
}

func (s *ConfigTestSuite) TestRunNewWithExistingConfigReturnsError() {
	//-- arrange
	err := data.SaveJSON(config.Config{}, s.ConfigLoader.ConfigPath)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-new")

	//-- assert
	s.RequireResultFail(fmt.Sprintf("config file \"%s\" already exists", s.ConfigLoader.ConfigPath))
}

func (s *ConfigTestSuite) TestRunNewCreatesNewConfig() {
	//-- arrange
	_ = os.Remove(s.ConfigLoader.ConfigPath)

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-new")

	//-- assert
	s.RequireResultPass()

	s.Assert().FileExists(s.ConfigLoader.ConfigPath)
	s.Assert().Contains(out.ReadLine(), "[+] Created New Config: "+s.ConfigLoader.ConfigPath)
}

func (s *ConfigTestSuite) TestRunInitVaultInvalidPathReturnsError() {
	//-- arrange
	CONFIG := config.Config{
		Version:   config.VERSION,
		VaultPath: file.NewPath(s.T(), ""),
	}
	err := data.SaveJSON(CONFIG, s.ConfigLoader.ConfigPath)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-init")

	//-- assert
	s.RequireResultFail("error initializing new vault")
}

func (s *ConfigTestSuite) TestRunInitVaultInitializesVault() {
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
	s.RunCommand("-init")

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), "[+] New Vault Initialized: "+CONFIG.VaultPath)
}
