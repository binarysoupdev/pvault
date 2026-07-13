package cmd_test

import (
	"fmt"
	"os"
	"pvault/cmd"
	"pvault/config"
	"pvault/vault"
	v1 "pvault/vault/database/version/v1"
	"pvault/vault/index"
	"pvault/vault/record"
	"regexp"
	"testing"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type VaultTestSuite struct {
	test.CommandSuite[*cmd.VaultCommand]
	ConfigLoader json.Loader[config.Config]
	Config       config.Config
}

func TestVaultCommandSuite(t *testing.T) {
	s := VaultTestSuite{
		ConfigLoader: json.NewLoader[config.Config](file.NewPath(t, "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewVaultCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *VaultTestSuite) SetupTest() {
	s.Config = config.Config{
		Version:    config.VERSION,
		VaultPath:  file.NewPath(s.T(), "vault"),
		BackupPath: file.NewPath(s.T(), ""),
	}

	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)
}

//=====================================

func (s *VaultTestSuite) TestRunFailErrorLoadingConfig() {
	//-- arrange
	err := os.Remove(s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("invalid config path")
}

func (s *VaultTestSuite) TestRunInitWithInvalidVaultPathReturnsError() {
	//-- arrange
	s.Config.VaultPath = file.NewPath(s.T(), "")
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-init")

	//-- assert
	s.RequireResultFail("error initializing new vault")
}

func (s *VaultTestSuite) TestRunInitValidInitializesVault() {
	//-- arrange
	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-init")

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), "[+] New Vault Initialized: "+s.Config.VaultPath)
	s.Assert().DirExists(s.Config.VaultPath)
}

func (s *VaultTestSuite) TestRunBackupFailsWithInvalidVault() {
	//-- arrange
	s.Config.VaultPath = file.NewPath(s.T(), "")
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-backup")

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *VaultTestSuite) TestRunBackupFailsWithInvalidBackupPath() {
	//-- arrange
	s.Config.BackupPath = file.CreateEmpty(s.T(), "backups.txt")
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	_, err = vault.InitializeNew(s.Config.VaultPath)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-backup")

	//-- assert
	s.RequireResultFail("error validating backup path")
}

func (s *VaultTestSuite) TestRunBackupPassesAndBacksUpVault() {
	//-- arrange
	DIR_REGEX := regexp.MustCompile(`"([^"]*)"`)

	_, err := vault.InitializeNew(s.Config.VaultPath)
	s.Require().NoError(err)

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-backup")

	//-- assert
	s.RequireResultPass()

	line := out.ReadLine()
	require.Contains(s.T(), line, "[+] Created Backup")

	match := DIR_REGEX.FindStringSubmatch(line)
	require.Len(s.T(), match, 2)
	assert.DirExists(s.T(), match[1])
}

func (s *VaultTestSuite) TestRunUpgradeFailsWithInvalidVault() {
	//-- arrange
	s.Config.VaultPath = file.NewPath(s.T(), "")
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-upgrade")

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *VaultTestSuite) TestRunUpgradeFailsWhenVaultIsUpToDate() {
	//-- arrange
	_, err := vault.InitializeNew(s.Config.VaultPath)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-upgrade")

	//-- assert
	s.RequireResultFail("vault is up-to-date")
}

func (s *VaultTestSuite) TestRunUpgradeFailsWithInvalidBackupPath() {
	//-- arrange
	s.Config.VaultPath = file.NewPath(s.T(), "")
	s.Config.BackupPath = file.CreateEmpty(s.T(), "backups.txt")

	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	err = v1.New(s.Config.VaultPath).Initialize(index.IndexMap{})
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-upgrade")

	//-- assert
	s.RequireResultFail("error validating backup path")
}

func (s *VaultTestSuite) TestRunUpgradePassesAndCreatesBackupAndUpgradesDatabase() {
	//-- arrange
	DIR_REGEX := regexp.MustCompile(`"([^"]*)"`)

	s.Config.VaultPath = file.NewPath(s.T(), "")
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	v1 := v1.New(s.Config.VaultPath)
	err = v1.Initialize(index.IndexMap{})
	s.Require().NoError(err)

	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.RunCommand("-upgrade")

	//-- assert
	s.RequireResultPass()

	line := out.ReadLine()
	require.Contains(s.T(), line, "[+] Created Backup")

	match := DIR_REGEX.FindStringSubmatch(line)
	require.Len(s.T(), match, 2)
	assert.DirExists(s.T(), match[1])

	v, err := vault.Open(s.Config.VaultPath)
	s.Require().NoError(err)

	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("[+] Vault Upgraded (@v%d -> @v%d)", v1.GetVersion(), v.Version()))
}

func (s *VaultTestSuite) TestRunValidateWithInvalidVaultFails() {
	//-- arrange
	s.Config.VaultPath = file.NewPath(s.T(), "")
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *VaultTestSuite) TestRunValidatePassPrintsVaultPathAndRecordCount() {
	//-- arrange
	v, err := vault.InitializeNew(s.Config.VaultPath)
	s.Require().NoError(err)

	rand := rand.New(0)

	NUM_RECORDS := 5
	for range NUM_RECORDS {
		err := v.SaveRecord(record.NewFromName(rand.ASCII(10)), rand.ASCII(30))
		s.Require().NoError(err)
	}

	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Vault verified at \"%s\" (@v%d)", v.Path, v.Version()))
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("[%d] records found", NUM_RECORDS))
}
