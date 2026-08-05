package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	cmd "pvault/app/commands/config"
	"pvault/app/config"
	"pvault/app/vault"
	"pvault/app/vault/database"
	db_v1 "pvault/app/vault/database/encoder/legacy/v1"
	"pvault/app/vault/index"
	"pvault/app/vault/meta"
	meta_v1 "pvault/app/vault/meta/encoder/v1"
	"pvault/util"
	"testing"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-commando/test"
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
		ConfigLoader: json.NewLoader[config.Config](filepath.Join(t.TempDir(), "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewConfigCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *ConfigTestSuite) SetupTest() {
	s.Config = config.Config{
		Version:    config.VERSION,
		VaultPath:  filepath.Join(s.T().TempDir(), "vault"),
		BackupPath: s.T().TempDir(),
		OutputPath: s.T().TempDir(),
	}
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))
}

//=====================================

func (s *ConfigTestSuite) TestRunNewFailsWithExistingConfig() {
	//-- act
	s.RunCommand("-new")

	//-- assert
	s.RequireResultFail(fmt.Sprintf("config file \"%s\" already exists", s.ConfigLoader.Path))
}

func (s *ConfigTestSuite) TestRunNewPassesAndCreatesNewConfig() {
	//-- arrange
	s.Require().NoError(os.Remove(s.ConfigLoader.Path))

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-new")

	//-- assert
	s.RequireResultPass()

	s.Assert().FileExists(s.ConfigLoader.Path)
	s.Assert().Contains(out.ReadLine(), "[+] Created New Config: "+s.ConfigLoader.Path)
}

func (s *ConfigTestSuite) TestRunFailsWhenConfigNotFound() {
	//-- arrange
	s.Require().NoError(os.Remove(s.ConfigLoader.Path))

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("invalid config path")
}

func (s *ConfigTestSuite) TestRunValidatePassesWithInvalidVaultAndPrintsError() {
	//-- arrange
	s.Config.VaultPath = s.T().TempDir()
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

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

func (s *ConfigTestSuite) TestRunValidatePassesWithOutOfDateVaultAndPrintsError() {
	//-- arrange
	PATH := s.T().TempDir()
	DATABASE := db_v1.Encoder{}

	META := meta.Metadata{
		DatabaseVersion: DATABASE.GetVersion(),
	}
	s.Require().NoError(meta.SaveMetadata(meta_v1.Encoder{}, PATH, META))
	s.Require().NoError(database.SaveIndex(DATABASE, PATH, index.IndexMap{}))

	s.Config.VaultPath = PATH
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	out := pipe.OpenStdout(3)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Loaded from \"%s\"", s.ConfigLoader.Path))

	vaultPath := out.ReadLine()
	s.Assert().Contains(vaultPath, s.Config.VaultPath)
	s.Assert().Contains(vaultPath, fmt.Sprintf("vault (@v%d) out-of-date", META.DatabaseVersion))
}

func (s *ConfigTestSuite) TestRunValidatePassesWithInvalidBackupPathAndPrintsError() {
	//-- arrange
	s.Config.BackupPath = filepath.Join(s.T().TempDir(), "backup.txt")
	s.Require().NoError(util.CreateEmptyFile(s.Config.BackupPath))
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

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

func (s *ConfigTestSuite) TestRunValidatePassesWithInvalidOutputPathAndPrintsError() {
	//-- arrange
	s.Config.OutputPath = "invalid"
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

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

func (s *ConfigTestSuite) TestRunValidatePassesAndPrintsConfig() {
	//-- arrange
	_, err := vault.InitializeNew(s.Config.VaultPath, "")
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
	s.Assert().Contains(vaultPath, fmt.Sprintf("verified (@v%d)", database.CURRENT_VERSION))

	backupPath := out.ReadLine()
	s.Assert().Contains(backupPath, s.Config.BackupPath)
	s.Assert().Contains(backupPath, "verified")

	outputPath := out.ReadLine()
	s.Assert().Contains(outputPath, s.Config.OutputPath)
	s.Assert().Contains(outputPath, "verified")
}
