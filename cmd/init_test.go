package cmd_test

import (
	"pvault/cmd"
	"pvault/config"
	"testing"

	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/suite"
)

type InitTestSuite struct {
	test.CommandSuite[*cmd.InitCommand]
}

func TestInitCommandSuite(t *testing.T) {
	suite.Run(t, &InitTestSuite{
		CommandSuite: test.NewCommandSuite(cmd.NewInitCommand()),
	})
}

//=====================================

func (s *InitTestSuite) TestRunInvalidVaultPath() {
	//-- arrange
	config.SetGlobal(config.Config{
		VaultPath: file.NewPath(s.T(), ""),
	})

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("error initializing new vault")
}

func (s *InitTestSuite) TestRunValid() {
	//-- arrange
	config.SetGlobal(config.Config{
		VaultPath: file.NewPath(s.T(), "vault"),
	})

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), "[+] New Vault Initialized")
}
