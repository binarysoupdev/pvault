package cmd_test

import (
	"os"
	"path/filepath"
	"pvault/cmd"
	"pvault/config"
	"pvault/data"
	"pvault/vault"
	"pvault/vault/record"
	"testing"

	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type UnlockTestSuite struct {
	test.CommandSuite[*cmd.UnlockCommand]
	Vault    vault.Vault
	Record   record.Record
	Password string
}

func TestUnlockCommandSuite(t *testing.T) {
	suite.Run(t, &UnlockTestSuite{
		CommandSuite: test.NewCommandSuite(cmd.NewUnlockCommand()),
	})
}

func (s *UnlockTestSuite) SetupTest() {
	config.SetGlobal(config.Config{
		VaultPath:  file.NewPath(s.T(), "vault"),
		OutputPath: file.NewPath(s.T(), ""),
	})

	rand := rand.New(0)
	s.Record = record.NewFromName(rand.ASCII(15))
	s.Password = rand.ASCII(30)

	var err error
	s.Vault, err = vault.InitializeNew(config.Global.VaultPath)
	s.Require().NoError(err)

	err = s.Vault.SaveRecord(s.Record, s.Password)
	s.Require().NoError(err)
}

//=====================================

func (s *UnlockTestSuite) TestRunInvalidVaultPath() {
	//-- arrange
	config.Global.VaultPath = "invalid"

	//-- act
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *UnlockTestSuite) TestRunInvalidNoResults() {
	//-- act
	s.RunCommand("-s", "no match")

	//-- assert
	s.RequireResultFail("no matches found")
}

func (s *UnlockTestSuite) TestRunVaultFileMissing() {
	//-- arrange
	err := os.Remove(s.Vault.RecordPath(s.Record.ID))
	s.Require().NoError(err)

	io := pipe.OpenStdio(1, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", s.Password)
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error reading record file")

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Enter PASSWORD")
}

func (s *UnlockTestSuite) TestRunIncorrectPassword() {
	//-- arrange
	io := pipe.OpenStdio(1, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", s.Password+"x")
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error decrypting ciphertext")

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Enter PASSWORD")
}

func (s *UnlockTestSuite) TestRunInvalidOutputPath() {
	//-- arrange
	config.Global.OutputPath = "invalid"

	io := pipe.OpenStdio(1, 3, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", s.Password)
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error creating output record")

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Enter PASSWORD")
	s.Assert().Contains(io.ReadLine(), "[=] Loaded Record: "+s.Record.ID.String())
}

func (s *UnlockTestSuite) TestRunValid() {
	//-- arrange
	OUTPUT_FILE := filepath.Join(config.Global.OutputPath, s.Record.ID.String()+".json")

	io := pipe.OpenStdio(1, 4, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", s.Password)
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultPass()

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Enter PASSWORD")
	s.Assert().Contains(io.ReadLine(), "[=] Loaded Record: "+s.Record.ID.String())
	s.Assert().Contains(io.ReadLine(), "[+] "+OUTPUT_FILE)

	record, err := data.LoadJSON[record.Record](OUTPUT_FILE)
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, record)
}
