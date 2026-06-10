package cmd_test

import (
	"path/filepath"
	"pvault/cfg"
	"pvault/cmd"
	"pvault/data"
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
	cfg.SetGlobal(cfg.Config{
		VaultPath:  file.NewPath(s.T(), ""),
		OutputPath: file.NewPath(s.T(), ""),
	})
}

func (s *CreateTestSuite) TestNameNotEmpty() {
	//-- act
	s.RunCommand("-name", "")

	//-- assert
	s.RequireResultFail("\"name\" cannot be empty")
}

func (s *CreateTestSuite) TestInvalidVaultPath() {
	//-- assert
	cfg.Global.VaultPath += "/invalid"

	//-- act
	s.RunCommand("-name", "foobar")

	//-- assert
	s.RequireResultFail("error creating vault record")
}

func (s *CreateTestSuite) TestInvalidOutputPath() {
	//-- assert
	cfg.Global.OutputPath += "/invalid"

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-name", "foobar")

	//-- assert
	s.RequireResultFail("error creating output record")
	s.Assert().Contains(out.ReadLine(), "[+] New Record: ")
}

func (s *CreateTestSuite) TestCreateRecord() {
	//-- arrange
	r := rand.New(0)
	NAME := r.ASCII(10)

	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.RunCommand("-name", NAME)

	//-- assert
	s.RequireResultPass()

	line := out.ReadLine()
	s.Require().Contains(line, "[+] New Record: ")

	ID := line[len(line)-36:]
	VAULT_FILE := filepath.Join(cfg.Global.VaultPath, ID+".json")
	OUTPUT_FILE := filepath.Join(cfg.Global.OutputPath, ID+".json")

	s.Assert().Contains(out.ReadLine(), "[+] "+OUTPUT_FILE)

	record, err := data.LoadJSON[vault.Record](VAULT_FILE)
	s.Require().NoError(err)
	s.Assert().Equal(ID, record.ID.String())
	s.Assert().Equal(NAME, record.Name)

	record, err = data.LoadJSON[vault.Record](OUTPUT_FILE)
	s.Require().NoError(err)
	s.Assert().Equal(ID, record.ID.String())
	s.Assert().Equal(NAME, record.Name)
}
