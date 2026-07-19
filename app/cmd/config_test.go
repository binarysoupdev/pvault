package cmd_test

import (
	"fmt"
	"os"
	"path/filepath"
	"pvault/app/cmd"
	v1 "pvault/app/vault/index/version1"
	"pvault/app/vault/local"
	"pvault/config"
	"testing"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	test.CommandSuite[*cmd.ConfigCommand]
	ConfigLoader json.Loader[config.Config]
	Config       config.Config
}

func TestConfigCommandSuite(t *testing.T) {
	s := ConfigTestSuite{
		ConfigLoader: json.NewLoader[config.Config](file.NewPath(t, "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewConfigCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *ConfigTestSuite) SetupTest() {
	s.Config = config.Config{
		Version:    config.VERSION,
		VaultPath:  file.NewPath(s.T(), "vault"),
		BackupPath: file.NewPath(s.T(), ""),
		OutputPath: file.NewPath(s.T(), ""),
	}

	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)
}

//=====================================

func (s *ConfigTestSuite) TestRunNewWithExistingConfigReturnsError() {
	//-- act
	s.RunCommand("-new")

	//-- assert
	s.RequireResultFail(fmt.Sprintf("config file \"%s\" already exists", s.ConfigLoader.Path))
}

func (s *ConfigTestSuite) TestRunNewCreatesNewConfig() {
	//-- arrange
	err := os.Remove(s.ConfigLoader.Path)
	s.Require().NoError(err)

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-new")

	//-- assert
	s.RequireResultPass()

	s.Assert().FileExists(s.ConfigLoader.Path)
	s.Assert().Contains(out.ReadLine(), "[+] Created New Config: "+s.ConfigLoader.Path)
}

func (s *ConfigTestSuite) TestRunNotNewConfigNotFoundReturnsError() {
	//-- arrange
	err := os.Remove(s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("invalid config path")
}

func (s *ConfigTestSuite) TestRunValidateConfigWithInvalidVaultPrintsError() {
	//-- arrange
	s.Config.VaultPath = file.NewPath(s.T(), "")
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	out := pipe.OpenStdout(3)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Loaded from \"%s\"", s.ConfigLoader.Path))

	vaultPath := out.ReadLine()
	s.Assert().Contains(vaultPath, s.Config.VaultPath)
	s.Assert().Contains(vaultPath, "error opening vault")
}

func (s *ConfigTestSuite) TestRunValidateConfigWithOutOfDateVaultPrintsError() {
	//-- arrange
	PATH := file.CreateEmpty(s.T(), v1.FILENAME)
	const LEGACY_VERSION = 1

	s.Config.VaultPath = filepath.Dir(PATH)
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	out := pipe.OpenStdout(3)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Loaded from \"%s\"", s.ConfigLoader.Path))

	vaultPath := out.ReadLine()
	s.Assert().Contains(vaultPath, s.Config.VaultPath)
	s.Assert().Contains(vaultPath, fmt.Sprintf("vault (@v%d) out-of-date", LEGACY_VERSION))
}

func (s *ConfigTestSuite) TestRunValidatePassWithInvalidBackupPathPrintsError() {
	//-- arrange
	s.Config.BackupPath = file.CreateEmpty(s.T(), "backup.txt")
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	out := pipe.OpenStdout(3)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Loaded from \"%s\"", s.ConfigLoader.Path))
	out.SkipLines(1)

	outputPath := out.ReadLine()
	s.Assert().Contains(outputPath, s.Config.BackupPath)
	s.Assert().Contains(outputPath, "path not a directory")
}

func (s *ConfigTestSuite) TestRunValidateWithInvalidOutputPathPrintsError() {
	//-- arrange
	s.Config.OutputPath = "invalid"
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	out := pipe.OpenStdout(4)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Loaded from \"%s\"", s.ConfigLoader.Path))
	out.SkipLines(2)

	outputPath := out.ReadLine()
	s.Assert().Contains(outputPath, s.Config.OutputPath)
	s.Assert().Contains(outputPath, "path not found")
}

func (s *ConfigTestSuite) TestRunValidatePassConfigPrintsConfig() {
	//-- arrange
	_, err := local.CreateNewVault(s.Config.VaultPath)
	s.Require().NoError(err)

	out := pipe.OpenStdout(4)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Loaded from \"%s\"", s.ConfigLoader.Path))

	vaultPath := out.ReadLine()
	s.Assert().Contains(vaultPath, s.Config.VaultPath)
	s.Assert().Contains(vaultPath, fmt.Sprintf("verified (@v%d)", local.CURRENT_VERSION))

	backupPath := out.ReadLine()
	s.Assert().Contains(backupPath, s.Config.BackupPath)
	s.Assert().Contains(backupPath, "verified")

	outputPath := out.ReadLine()
	s.Assert().Contains(outputPath, s.Config.OutputPath)
	s.Assert().Contains(outputPath, "verified")
}
