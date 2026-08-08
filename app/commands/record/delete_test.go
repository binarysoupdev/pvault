package record_test

import (
	"fmt"
	"os"
	cmd "pvault/app/commands/record"
	"pvault/app/config"
	"pvault/app/vault"
	"pvault/app/vault/database"
	db_v1 "pvault/app/vault/database/encoder/legacy/v1"
	"pvault/app/vault/index"
	"pvault/app/vault/meta"
	meta_v1 "pvault/app/vault/meta/encoder/v1"
	record_v2 "pvault/app/vault/record/record/v2"
	"testing"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type DeleteTestSuite struct {
	test.CommandSuite[*cmd.DeleteCommand]
	ConfigLoader json.Loader[config.Config]
	Config       config.Config

	Vault  vault.Vault
	Record record_v2.Record
}

func TestDeleteCommandSuite(t *testing.T) {
	s := DeleteTestSuite{
		ConfigLoader: json.NewLoader[config.Config](file.NewPath(t, "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewDeleteCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *DeleteTestSuite) SetupTest() {
	s.Config = config.Config{
		Version:    config.VERSION,
		VaultPath:  file.NewPath(s.T(), "vault"),
		OutputPath: file.NewPath(s.T(), ""),
	}
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	rand := rand.New(0)
	s.Record = record_v2.NewEmptyRecord(rand.ASCII(15))

	s.Vault, err = vault.InitializeNew(s.Config.VaultPath, "")
	s.Require().NoError(err)

	err = s.Vault.SaveRecord(s.Record, rand.ASCII(30))
	s.Require().NoError(err)
}

//=====================================

func (s *DeleteTestSuite) TestRunFailsWhenErrorLoadingConfig() {
	//-- arrange
	err := os.Remove(s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("invalid config path")
}

func (s *DeleteTestSuite) TestRunFailsWithInvalidVaultPath() {
	//-- arrange
	s.Config.VaultPath = "invalid"
	err := json.MarshalFile(s.Config, s.ConfigLoader.Path)
	s.Require().NoError(err)

	//-- act
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *DeleteTestSuite) TestRunFailsWhenVaultOutOfDate() {
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
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail(fmt.Sprintf("vault (@v%d) out-of-date", DATABASE.GetVersion()))
}

func (s *DeleteTestSuite) TestRunFailsWhenNoResults() {
	//-- act
	s.RunCommand("-s", "no match")

	//-- assert
	s.RequireResultFail("no matches found")
}

func (s *DeleteTestSuite) TestRunFailsWithIncorrectConfirmName() {
	//-- arrange
	io := pipe.OpenStdio(1, 2, true)
	defer io.Close()

	//-- act
	io.Queue("NAME: ", s.Record.Name+"x")
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("names do not match")

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Confirm NAME: "+s.Record.Name+"x")
}

func (s *DeleteTestSuite) TestRunPassesAndDeletesRecord() {
	//-- arrange
	io := pipe.OpenStdio(1, 3, true)
	defer io.Close()

	//-- act
	io.Queue("NAME: ", s.Record.Name)
	io.EndQueue()

	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultPass()

	s.Assert().Contains(io.ReadLine(), s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "Confirm NAME: "+s.Record.Name)
	s.Assert().Contains(io.ReadLine(), "[-] Deleted Record: "+s.Record.ID.String())

	err := s.Vault.LoadIndex()
	s.Require().NoError(err)

	_, err = s.Vault.LoadRecord(s.Record.Name, "")
	s.Assert().ErrorContains(err, fmt.Sprintf("name \"%s\" not found", s.Record.Name))
}
