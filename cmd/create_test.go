package cmd_test

import (
	"path/filepath"
	"pvault/cmd"
	"pvault/config"
	"pvault/vault"
	"testing"

	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type CreateTestSuite struct {
	test.CommandSuite[*cmd.CreateCommand]
}

func TestCreateCommandSuite(t *testing.T) {
	suite.Run(t, &CreateTestSuite{
		CommandSuite: test.NewCommandSuite(cmd.NewCreateCommand()),
	})
}

func (s *CreateTestSuite) SetupTest() {
	config.SetGlobal(config.Config{
		VaultPath:  file.NewPath(s.T(), "vault"),
		OutputPath: file.NewPath(s.T(), ""),
	})

	err := vault.InitializeNew(config.Global.VaultPath)
	s.Require().NoError(err)
}

//=====================================

func (s *CreateTestSuite) TestRunNameNotEmpty() {
	//-- act
	s.RunCommand("-name", "")

	//-- assert
	s.RequireResultFail("\"name\" cannot be empty")
}

func (s *CreateTestSuite) TestRunInvalidVaultPath() {
	//-- arrange
	rand := rand.New(0)
	NAME := rand.ASCII(15)

	config.Global.VaultPath = "invalid"

	//-- act
	s.RunCommand("-name", NAME)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *CreateTestSuite) TestRunInvalidOutputPath() {
	//-- arrange
	rand := rand.New(0)
	NAME := rand.ASCII(15)

	config.Global.OutputPath = "invalid"

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-name", NAME)

	//-- assert
	s.RequireResultFail("error creating output record")
	s.Assert().Contains(out.ReadLine(), "[+] New Record: ")
}

func (s *CreateTestSuite) TestRunValid() {
	//-- arrange
	rand := rand.New(0)
	NAME := rand.ASCII(15)

	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.RunCommand("-name", NAME)

	//-- assert
	s.RequireResultPass()

	line := out.ReadLine()
	s.Require().Contains(line, "[+] New Record: ")

	ID := line[len(line)-36:]
	VAULT_FILE := filepath.Join(config.Global.VaultPath, ID+".json")
	OUTPUT_FILE := filepath.Join(config.Global.OutputPath, ID+".json")

	s.Assert().Contains(out.ReadLine(), "[+] "+OUTPUT_FILE)
	s.Assert().FileExists(VAULT_FILE)
	s.Assert().FileExists(OUTPUT_FILE)
}

func (s *CreateTestSuite) TestRunExistingNameInvalid() {
	//-- arrange
	rand := rand.New(0)
	NAME := rand.ASCII(15)

	v, err := vault.Open(config.Global.VaultPath)
	s.Require().NoError(err)

	v.SaveRecord(vault.EmptyRecord(NAME))
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-name", NAME)

	//-- assert
	s.RequireResultFail("error saving vault record")
}
