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

	_, err := vault.InitializeNew(config.Global.VaultPath)
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

func (s *CreateTestSuite) TestRunIncorrectVerifyPassword() {
	//-- arrange
	rand := rand.New(0)
	NAME := rand.ASCII(15)
	PASSWORD := rand.ASCII(30)

	io := pipe.OpenStdio(2, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.Queue("PASSWORD: ", PASSWORD+"x")
	io.EndQueue()

	s.RunCommand("-name", NAME)

	//-- assert
	s.RequireResultFail("passwords do not match")
	s.Assert().Contains(io.ReadLine(), "New PASSWORD")
	s.Assert().Contains(io.ReadLine(), "Verify PASSWORD")
}

func (s *CreateTestSuite) TestRunInvalidOutputPath() {
	//-- arrange
	rand := rand.New(0)
	NAME := rand.ASCII(15)
	PASSWORD := rand.ASCII(30)

	config.Global.OutputPath = "invalid"

	io := pipe.OpenStdio(2, 3, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.Queue("PASSWORD: ", PASSWORD)
	io.EndQueue()

	s.RunCommand("-name", NAME)

	//-- assert
	s.RequireResultFail("error creating output record")
	io.SkipLines(2)
	s.Assert().Contains(io.ReadLine(), "[+] New Record: ")
}

func (s *CreateTestSuite) TestRunValid() {
	//-- arrange
	rand := rand.New(0)
	NAME := rand.ASCII(15)
	PASSWORD := rand.ASCII(30)

	io := pipe.OpenStdio(2, 4, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.Queue("PASSWORD: ", PASSWORD)
	io.EndQueue()

	s.RunCommand("-name", NAME)

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(io.ReadLine(), "New PASSWORD")
	s.Assert().Contains(io.ReadLine(), "Verify PASSWORD")

	line := io.ReadLine()
	s.Require().Contains(line, "[+] New Record: ")

	ID := line[len(line)-36:]
	VAULT_FILE := filepath.Join(config.Global.VaultPath, ID+".json")
	OUTPUT_FILE := filepath.Join(config.Global.OutputPath, ID+".json")

	s.Assert().Contains(io.ReadLine(), "[+] "+OUTPUT_FILE)
	s.Assert().FileExists(VAULT_FILE)
	s.Assert().FileExists(OUTPUT_FILE)
}
