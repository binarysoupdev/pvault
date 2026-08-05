package vault_test

import (
	"fmt"
	"os"
	"path/filepath"
	cmd "pvault/app/commands/vault"
	"pvault/app/config"
	"pvault/app/vault"
	"pvault/app/vault/database"
	db_v1 "pvault/app/vault/database/encoder/legacy/v1"
	"pvault/app/vault/index"
	record_v2 "pvault/app/vault/record/record/v2"
	"pvault/util"
	"regexp"
	"testing"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/pipe"
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
		ConfigLoader: json.NewLoader[config.Config](filepath.Join(t.TempDir(), "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewVaultCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *VaultTestSuite) SetupTest() {
	s.Config = config.Config{
		Version:    config.VERSION,
		VaultPath:  filepath.Join(s.T().TempDir(), "vault"),
		BackupPath: s.T().TempDir(),
	}

	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)
}

//=====================================

func (s *VaultTestSuite) TestRunFailsWhenErrorLoadingConfig() {
	//-- arrange
	err := os.Remove(s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("invalid config path")
}

func (s *VaultTestSuite) TestRunInitFailsWithInvalidVaultPath() {
	//-- arrange
	s.Config.VaultPath = s.T().TempDir()
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-init")

	//-- assert
	s.RequireResultFail("error initializing new vault")
}

func (s *VaultTestSuite) TestRunInitPassesAndInitializesVault() {
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
	s.Config.VaultPath = s.T().TempDir()
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-backup")

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *VaultTestSuite) TestRunBackupFailsWithInvalidBackupPath() {
	//-- arrange
	s.Config.BackupPath = filepath.Join(s.T().TempDir(), "backups.txt")
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))
	s.Require().NoError(util.CreateEmptyFile(s.Config.BackupPath))

	_, err := vault.InitializeNew(s.Config.VaultPath, "")
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-backup")

	//-- assert
	s.RequireResultFail("error validating \"config.backup_path\"")
}

func (s *VaultTestSuite) TestRunBackupPassesAndBacksUpVault() {
	//-- arrange
	DIR_REGEX := regexp.MustCompile(`"([^"]*)"`)

	_, err := vault.InitializeNew(s.Config.VaultPath, "")
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
	s.Config.VaultPath = s.T().TempDir()
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-upgrade")

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *VaultTestSuite) TestRunUpgradeFailsWhenVaultIsUpToDate() {
	//-- arrange
	_, err := vault.InitializeNew(s.Config.VaultPath, "")
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-upgrade")

	//-- assert
	s.RequireResultFail("vault is up-to-date")
}

func (s *VaultTestSuite) TestRunUpgradeFailsWithInvalidBackupPath() {
	//-- arrange
	s.Config.VaultPath = s.T().TempDir()
	s.Config.BackupPath = filepath.Join(s.T().TempDir(), "backups.txt")
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))
	s.Require().NoError(util.CreateEmptyFile(s.Config.BackupPath))

	s.Require().NoError(database.SaveIndex(db_v1.Encoder{}, s.Config.VaultPath, index.IndexMap{}))

	//-- act
	s.RunCommand("-upgrade")

	//-- assert
	s.RequireResultFail("error validating \"config.backup_path\"")
}

func (s *VaultTestSuite) TestRunUpgradePassesAndCreatesBackupAndUpgradesDatabase() {
	//-- arrange
	DIR_REGEX := regexp.MustCompile(`"([^"]*)"`)

	s.Config.VaultPath = s.T().TempDir()
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	db := db_v1.Encoder{}
	s.Require().NoError(database.SaveIndex(db, s.Config.VaultPath, index.IndexMap{}))

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

	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("[+] Vault Upgraded (@v%d -> @v%d)", db.GetVersion(), v.GetVersion()))
}

func (s *VaultTestSuite) TestRunValidateFailsWithInvalidVault() {
	//-- arrange
	s.Config.VaultPath = s.T().TempDir()
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *VaultTestSuite) TestRunValidatePassesAndPrintsVaultPathAndRecordCount() {
	//-- arrange
	v, err := vault.InitializeNew(s.Config.VaultPath, "")
	s.Require().NoError(err)

	const NUM_RECORDS = 5
	const PASSWORD = "Password123!"

	for i := range NUM_RECORDS {
		err := v.SaveRecord(record_v2.NewEmptyRecord(fmt.Sprintf("name_%d", i)), PASSWORD)
		s.Require().NoError(err)
	}

	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Vault verified at \"%s\" (@v%d)", v.Path, v.GetVersion()))
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("[%d] records found", NUM_RECORDS))
}
