package cmd_test

import (
	"pvault/cmd"
	"pvault/config"
	"pvault/data"
	"testing"

	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
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
