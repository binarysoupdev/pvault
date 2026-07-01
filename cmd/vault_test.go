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
	"pvault/vault/record"
	"testing"

	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type VaultTestSuite struct {
	test.CommandSuite[*cmd.VaultCommand]
	ConfigLoader config.Loader[config.Config]
	Config       config.Config
}

func TestVaultCommandSuite(t *testing.T) {
	s := VaultTestSuite{
		ConfigLoader: config.NewLoader[config.Config](file.NewPath(t, "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewVaultCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *VaultTestSuite) SetupTest() {
	s.Config = config.Config{
		Version:   config.VERSION,
		VaultPath: file.NewPath(s.T(), "vault"),
	}

	err := data.SaveJSON(s.Config, s.ConfigLoader.ConfigPath)
	s.Require().NoError(err)
}

//=====================================

func (s *VaultTestSuite) TestRunInitWithInvalidVaultPathReturnsError() {
	//-- arrange
	s.Config.VaultPath = file.NewPath(s.T(), "")
	err := data.SaveJSON(s.Config, s.ConfigLoader.ConfigPath)
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

func (s *VaultTestSuite) TestRunValidateWithInvalidVaultFails() {
	//-- arrange
	s.Config.VaultPath = file.NewPath(s.T(), "")
	err := data.SaveJSON(s.Config, s.ConfigLoader.ConfigPath)
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

func (s *VaultTestSuite) TestRunUpgradeWithInvalidVaultFails() {
	//-- arrange
	s.Config.VaultPath = file.NewPath(s.T(), "")
	err := data.SaveJSON(s.Config, s.ConfigLoader.ConfigPath)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-upgrade")

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *VaultTestSuite) TestRunUpgradeWhereVaultIsUpToDateFails() {
	//-- arrange
	_, err := vault.InitializeNew(s.Config.VaultPath)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-upgrade")

	//-- assert
	s.RequireResultFail("vault is up-to-date")
}

func (s *VaultTestSuite) TestRunUpgradePassCreatesBackupAndUpgradesDatabase() {
	//-- arrange
	err := os.Mkdir(s.Config.VaultPath, 0755)
	s.Require().NoError(err)

	v1 := version1.NewDatabase(s.Config.VaultPath)

	file, err := os.Create(v1.IndexPath())
	s.Require().NoError(err)
	file.Close()

	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.RunCommand("-upgrade")

	//-- assert
	s.RequireResultPass()

	backup := filepath.Join(s.Config.VaultPath, fmt.Sprintf("version_%d", v1.GetVersion()))
	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("[+] Created Backup \"%s\"", backup))
	s.Assert().DirExists(backup)

	v, err := vault.Open(s.Config.VaultPath)
	s.Require().NoError(err)

	s.Assert().Contains(out.ReadLine(), fmt.Sprintf("[+] Vault Upgraded (@v%d -> @v%d)", v1.GetVersion(), v.Version()))
}
