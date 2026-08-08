package cmds_test

import (
	"fmt"
	"os"
	"path/filepath"
	"pvault/app/cmds"
	"pvault/app/config"
	"pvault/app/vault"
	"pvault/app/vault/database"
	db_v1 "pvault/app/vault/database/encoder/legacy/v1"
	"pvault/app/vault/index"
	"pvault/app/vault/meta"
	meta_v1 "pvault/app/vault/meta/encoder/v1"
	record_v2 "pvault/app/vault/record/record/v2"
	"pvault/util"
	"regexp"
	"testing"
	"time"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type VaultTestSuite struct {
	test.CommandSuite[*cmds.VaultCommand]
	ConfigLoader json.Loader[config.Config]
	Config       config.Config
	Dir          string
}

func TestVaultCommandSuite(t *testing.T) {
	s := VaultTestSuite{
		ConfigLoader: json.NewLoader[config.Config](filepath.Join(t.TempDir(), "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmds.NewVaultCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *VaultTestSuite) SetupTest() {
	s.Dir = "vault"
	s.Config = config.Config{
		Version:    config.VERSION,
		VaultPath:  filepath.Join(s.T().TempDir(), s.Dir),
		BackupPath: s.T().TempDir(),
	}

	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))
}

//=====================================

func (s *VaultTestSuite) TestRunFailsWhenConfigNotFound() {
	//-- arrange
	s.Require().NoError(os.Remove(s.ConfigLoader.Path))

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("invalid config path")
}

func (s *VaultTestSuite) TestRunFailsWhenConfigVersionUnsupported() {
	//-- arrange
	s.Config.Version = config.VERSION + 1
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail(fmt.Sprintf("unsupported config version \"%d\"", s.Config.Version))
}

func (s *VaultTestSuite) TestRunSetNicknameFailsWithInvalidVault() {
	//-- arrange
	const NICKNAME = "nickname"

	s.Config.VaultPath = s.T().TempDir()
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	//-- act
	s.RunCommand("-nickname", NICKNAME)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *VaultTestSuite) TestRunSetNicknamePassesAndSetsNickname() {
	//-- arrange
	const NICKNAME = "nickname"

	_, err := vault.InitializeNew(s.Config.VaultPath, "")
	s.Require().NoError(err)

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-nickname", NICKNAME)

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("[+] Set Nickname: %s", NICKNAME))

	v, err := vault.Open(s.Config.VaultPath)
	s.Require().NoError(err)

	s.Assert().Equal(NICKNAME, v.Meta.Nickname)
}

func (s *VaultTestSuite) TestRunInitFailsWithInvalidVaultPath() {
	//-- arrange
	s.Config.VaultPath = s.T().TempDir()
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	//-- act
	s.RunCommand("-init")

	//-- assert
	s.RequireResultFail("error initializing new vault")
}

func (s *VaultTestSuite) TestRunInitPassesAndInitializesVaultWithDefaultNickname() {
	//-- arrange
	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-init")

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("[+] New Vault \"%s\" Initialized: %s", s.Dir, s.Config.VaultPath))
	s.Assert().DirExists(s.Config.VaultPath)
}

func (s *VaultTestSuite) TestRunInitPassesAndInitializesVaultWithNickname() {
	//-- arrange
	const NICKNAME = "nickname"

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-init", "-nickname", NICKNAME)

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("[+] New Vault \"%s\" Initialized: %s", NICKNAME, s.Config.VaultPath))
	s.Assert().DirExists(s.Config.VaultPath)
}

func (s *VaultTestSuite) TestRunBackupFailsWithInvalidVault() {
	//-- arrange
	s.Config.VaultPath = s.T().TempDir()
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

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
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

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
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

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
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *VaultTestSuite) TestRunValidateFailsWhenVaultOutOfDate() {
	//-- arrange
	s.Config.VaultPath = s.T().TempDir()
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	DATABASE := db_v1.Encoder{}

	META := meta.Metadata{
		DatabaseVersion: DATABASE.GetVersion(),
	}
	s.Require().NoError(meta.SaveMetadata(meta_v1.Encoder{}, s.Config.VaultPath, META))
	s.Require().NoError(database.SaveIndex(DATABASE, s.Config.VaultPath, index.IndexMap{}))

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail(fmt.Sprintf("vault (@v%d) out-of-date", DATABASE.GetVersion()))
}

func (s *VaultTestSuite) TestRunValidatePassesAndPrintsVaultPathAndRecordCount() {
	//-- arrange
	const NUM_RECORDS = 5
	const PASSWORD = "Password123!"
	const NICKNAME = "nickname"

	v, err := vault.InitializeNew(s.Config.VaultPath, NICKNAME)
	s.Require().NoError(err)

	for i := range NUM_RECORDS {
		s.Require().NoError(v.SaveRecord(record_v2.NewEmptyRecord(fmt.Sprintf("name_%d", i)), PASSWORD))
	}

	out := pipe.OpenStdout(3)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Vault \"%s\" verified at \"%s\" (@v%d)", NICKNAME, v.Path, v.GetVersion()))
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Created on %s", v.Meta.CreationDate.Format(time.DateOnly)))
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("[%d] records found", NUM_RECORDS))
}
