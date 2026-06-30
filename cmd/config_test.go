package cmd_test

import (
	"fmt"
	"os"
	"path/filepath"
	"pvault/cmd"
	"pvault/config"
	"pvault/data"
	"pvault/vault"
	"pvault/vault/data/version1"
	"testing"

	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	test.CommandSuite[*cmd.ConfigCommand]
	ConfigLoader config.Loader[config.Config]
	Config       config.Config
}

func TestConfigCommandSuite(t *testing.T) {
	s := ConfigTestSuite{
		ConfigLoader: config.NewLoader[config.Config](file.NewPath(t, "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewConfigCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *ConfigTestSuite) SetupTest() {
	s.Config = config.Config{
		Version:    config.VERSION,
		VaultPath:  file.NewPath(s.T(), "vault"),
		OutputPath: file.NewPath(s.T(), ""),
	}

	err := data.SaveJSON(s.Config, s.ConfigLoader.ConfigPath)
	s.Require().NoError(err)
}

//=====================================

func (s *ConfigTestSuite) TestRunNewWithExistingConfigReturnsError() {
	//-- act
	s.RunCommand("-new")

	//-- assert
	s.RequireResultFail(fmt.Sprintf("config file \"%s\" already exists", s.ConfigLoader.ConfigPath))
}

func (s *ConfigTestSuite) TestRunNewCreatesNewConfig() {
	//-- arrange
	err := os.Remove(s.ConfigLoader.ConfigPath)
	s.Require().NoError(err)

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-new")

	//-- assert
	s.RequireResultPass()

	s.Assert().FileExists(s.ConfigLoader.ConfigPath)
	s.Assert().Contains(out.ReadLine(), "[+] Created New Config: "+s.ConfigLoader.ConfigPath)
}

func (s *ConfigTestSuite) TestRunNotNewConfigNotFoundReturnsError() {
	//-- arrange
	err := os.Remove(s.ConfigLoader.ConfigPath)
	s.Require().NoError(err)

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("error loading config")
}

func (s *ConfigTestSuite) TestRunValidateConfigWithInvalidVaultPrintsError() {
	//-- arrange
	s.Config.VaultPath = file.NewPath(s.T(), "")
	err := data.SaveJSON(s.Config, s.ConfigLoader.ConfigPath)
	s.Require().NoError(err)

	out := pipe.OpenStdout(3)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Loaded from \"%s\"", s.ConfigLoader.ConfigPath))

	vaultPath := out.ReadLine()
	s.Assert().Contains(vaultPath, s.Config.VaultPath)
	s.Assert().Contains(vaultPath, "error opening vault")
}

func (s *ConfigTestSuite) TestRunValidateConfigWithOutOfDateVaultPrintsError() {
	//-- arrange
	PATH := file.CreateEmpty(s.T(), version1.INDEX_FILE)
	const LEGACY_VERSION = 1

	s.Config.VaultPath = filepath.Dir(PATH)
	err := data.SaveJSON(s.Config, s.ConfigLoader.ConfigPath)
	s.Require().NoError(err)

	out := pipe.OpenStdout(3)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Loaded from \"%s\"", s.ConfigLoader.ConfigPath))

	vaultPath := out.ReadLine()
	s.Assert().Contains(vaultPath, s.Config.VaultPath)
	s.Assert().Contains(vaultPath, fmt.Sprintf("vault (@v%d) out-of-date", LEGACY_VERSION))
}

func (s *ConfigTestSuite) TestRunValidateWithInvalidOutputPathPrintsError() {
	//-- arrange
	s.Config.OutputPath = "invalid"
	err := data.SaveJSON(s.Config, s.ConfigLoader.ConfigPath)
	s.Require().NoError(err)

	out := pipe.OpenStdout(3)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Loaded from \"%s\"", s.ConfigLoader.ConfigPath))
	out.SkipLines(1)

	outputPath := out.ReadLine()
	s.Assert().Contains(outputPath, s.Config.OutputPath)
	s.Assert().Contains(outputPath, "path not found")
}

func (s *ConfigTestSuite) TestRunValidateConfigPrintsConfig() {
	//-- arrange
	_, err := vault.InitializeNew(s.Config.VaultPath)
	s.Require().NoError(err)

	out := pipe.OpenStdout(3)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Loaded from \"%s\"", s.ConfigLoader.ConfigPath))

	vaultPath := out.ReadLine()
	s.Assert().Contains(vaultPath, s.Config.VaultPath)
	s.Assert().Contains(vaultPath, fmt.Sprintf("verified (@v%d)", vault.CURRENT_VERSION))

	outputPath := out.ReadLine()
	s.Assert().Contains(outputPath, s.Config.OutputPath)
	s.Assert().Contains(outputPath, "verified")
}
