package cmds_test

import (
	"fmt"
	"os"
	"path/filepath"
	"pvault/app/cmds"
	"pvault/app/config"

	"pvault/vault"
	"pvault/vault/database"
	db_v1 "pvault/vault/database/encoder/legacy/v1"
	"pvault/vault/index"
	"pvault/vault/meta"
	meta_v1 "pvault/vault/meta/encoder/v1"
	record_v2 "pvault/vault/record/record/v2"
	"regexp"
	"testing"
	"time"

	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/go-extensions/file"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/tinsel/pipe"
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

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-nickname", NICKNAME)

	//-- assert
	s.RequireResultFail("vault not found")
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Vault Path: \"%s\"", s.Config.VaultPath))
}

func (s *VaultTestSuite) TestRunSetNicknamePassesAndSetsNickname() {
	//-- arrange
	const NICKNAME = "nickname"

	_, err := vault.InitializeNew(s.Config.VaultPath, "")
	s.Require().NoError(err)

	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.RunCommand("-nickname", NICKNAME)

	//-- assert
	s.RequireResultPass()

	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Vault Path: \"%s\"", s.Config.VaultPath))
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("[+] Set Nickname: %s", NICKNAME))

	v, err := vault.Open(s.Config.VaultPath)
	s.Require().NoError(err)

	s.Assert().Equal(NICKNAME, v.Meta.Nickname)
}

func (s *VaultTestSuite) TestRunInitFailsWithInvalidVaultPath() {
	//-- arrange
	s.Config.VaultPath = s.T().TempDir()
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-init")

	//-- assert
	s.RequireResultFail("error initializing new vault")
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Vault Path: \"%s\"", s.Config.VaultPath))
}

func (s *VaultTestSuite) TestRunInitPassesAndInitializesVaultWithDefaultNickname() {
	//-- arrange
	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.RunCommand("-init")

	//-- assert
	s.RequireResultPass()

	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Vault Path: \"%s\"", s.Config.VaultPath))
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("[+] New Vault \"%s\" Initialized: %s", s.Dir, s.Config.VaultPath))

	s.Assert().DirExists(s.Config.VaultPath)
}

func (s *VaultTestSuite) TestRunInitPassesAndInitializesVaultWithNickname() {
	//-- arrange
	const NICKNAME = "nickname"

	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.RunCommand("-init", "-nickname", NICKNAME)

	//-- assert
	s.RequireResultPass()

	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Vault Path: \"%s\"", s.Config.VaultPath))
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("[+] New Vault \"%s\" Initialized: %s", NICKNAME, s.Config.VaultPath))

	s.Assert().DirExists(s.Config.VaultPath)
}

func (s *VaultTestSuite) TestRunBackupFailsWithInvalidVault() {
	//-- arrange
	s.Config.VaultPath = s.T().TempDir()
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-backup")

	//-- assert
	s.RequireResultFail("vault not found")
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Vault Path: \"%s\"", s.Config.VaultPath))
}

func (s *VaultTestSuite) TestRunBackupFailsWithInvalidBackupPath() {
	//-- arrange
	s.Config.BackupPath = filepath.Join(s.T().TempDir(), "backups.txt")
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))
	s.Require().NoError(file.CreateEmpty(s.Config.BackupPath))

	out := pipe.OpenStdout(1)
	defer out.Close()

	_, err := vault.InitializeNew(s.Config.VaultPath, "")
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-backup")

	//-- assert
	s.RequireResultFail("error validating \"config.backup_path\"")
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Vault Path: \"%s\"", s.Config.VaultPath))
}

func (s *VaultTestSuite) TestRunBackupPassesAndBacksUpVault() {
	//-- arrange
	DIR_REGEX := regexp.MustCompile(`"([^"]*)"`)

	_, err := vault.InitializeNew(s.Config.VaultPath, "")
	s.Require().NoError(err)

	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.RunCommand("-backup")

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Vault Path: \"%s\"", s.Config.VaultPath))

	line := out.ReadLine()
	s.Require().Contains(line, "[+] Created Backup")

	match := DIR_REGEX.FindStringSubmatch(line)
	s.Require().Len(match, 2)
	s.Assert().DirExists(match[1])
}

func (s *VaultTestSuite) TestRunUpgradeFailsWithInvalidVault() {
	//-- arrange
	s.Config.VaultPath = s.T().TempDir()
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-upgrade")

	//-- assert
	s.RequireResultFail("vault not found")
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Vault Path: \"%s\"", s.Config.VaultPath))
}

func (s *VaultTestSuite) TestRunUpgradeFailsWhenVaultIsUpToDate() {
	//-- arrange
	_, err := vault.InitializeNew(s.Config.VaultPath, "")
	s.Require().NoError(err)

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-upgrade")

	//-- assert
	s.RequireResultFail("vault is up-to-date")
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Vault Path: \"%s\"", s.Config.VaultPath))
}

func (s *VaultTestSuite) TestRunUpgradeFailsWithInvalidBackupPath() {
	//-- arrange
	s.Config.VaultPath = s.T().TempDir()
	s.Config.BackupPath = filepath.Join(s.T().TempDir(), "backups.txt")
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))
	s.Require().NoError(file.CreateEmpty(s.Config.BackupPath))

	s.Require().NoError(database.SaveIndex(db_v1.Encoder{}, s.Config.VaultPath, index.IndexMap{}))

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-upgrade")

	//-- assert
	s.RequireResultFail("error validating \"config.backup_path\"")
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Vault Path: \"%s\"", s.Config.VaultPath))
}

func (s *VaultTestSuite) TestRunUpgradePassesAndCreatesBackupAndUpgradesDatabase() {
	//-- arrange
	DIR_REGEX := regexp.MustCompile(`"([^"]*)"`)

	s.Config.VaultPath = s.T().TempDir()
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	db := db_v1.Encoder{}
	s.Require().NoError(database.SaveIndex(db, s.Config.VaultPath, index.IndexMap{}))

	out := pipe.OpenStdout(3)
	defer out.Close()

	//-- act
	s.RunCommand("-upgrade")

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Vault Path: \"%s\"", s.Config.VaultPath))

	line := out.ReadLine()
	s.Require().Contains(line, "[+] Created Backup")

	match := DIR_REGEX.FindStringSubmatch(line)
	s.Require().Len(match, 2)
	s.Assert().DirExists(match[1])

	v, err := vault.Open(s.Config.VaultPath)
	s.Require().NoError(err)

	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("[+] Vault Upgraded (@%d -> @%d)", db.GetVersion(), v.GetVersion()))
}

func (s *VaultTestSuite) TestRunValidateFailsWithInvalidVault() {
	//-- arrange
	s.Config.VaultPath = s.T().TempDir()
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("vault not found")
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Vault Path: \"%s\"", s.Config.VaultPath))
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

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail(fmt.Sprintf("vault (@v%d) out-of-date", DATABASE.GetVersion()))
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Vault Path: \"%s\"", s.Config.VaultPath))
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

	out := pipe.OpenStdout(4)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()

	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Vault Path: \"%s\"", s.Config.VaultPath))
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Verified vault \"%s\" (@v%d)", NICKNAME, v.GetVersion()))
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("Created on %s", v.Meta.CreationDate.Format(time.DateOnly)))
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("[%d] records found", NUM_RECORDS))
}
