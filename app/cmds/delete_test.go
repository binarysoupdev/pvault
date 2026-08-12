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
	"testing"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/suite"
)

type DeleteTestSuite struct {
	test.CommandSuite[*cmds.DeleteCommand]
	ConfigLoader json.Loader[config.Config]
	Config       config.Config

	Vault  vault.Vault
	Record record_v2.Record
}

func TestDeleteCommandSuite(t *testing.T) {
	s := DeleteTestSuite{
		ConfigLoader: json.NewLoader[config.Config](filepath.Join(t.TempDir(), "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmds.NewDeleteCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *DeleteTestSuite) SetupTest() {
	s.Config = config.Config{
		Version:    config.VERSION,
		VaultPath:  filepath.Join(s.T().TempDir(), "vault"),
		OutputPath: s.T().TempDir(),
	}
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	s.Record = record_v2.NewEmptyRecord("name")

	var err error
	s.Vault, err = vault.InitializeNew(s.Config.VaultPath, "")
	s.Require().NoError(err)

	s.Require().NoError(s.Vault.SaveRecord(s.Record, "Password123!"))
}

//=====================================

func (s *DeleteTestSuite) TestRunFailsWhenConfigNotFound() {
	//-- arrange
	s.Require().NoError(os.Remove(s.ConfigLoader.Path))

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("invalid config path")
}

func (s *DeleteTestSuite) TestRunFailsWhenConfigVersionUnsupported() {
	//-- arrange
	s.Config.Version = config.VERSION + 1
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail(fmt.Sprintf("unsupported config version \"%d\"", s.Config.Version))
}

func (s *DeleteTestSuite) TestRunFailsWithInvalidVaultPath() {
	//-- arrange
	s.Config.VaultPath = "invalid"
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	//-- act
	s.RunCommand("-s", s.Record.Name)

	//-- assert
	s.RequireResultFail("vault not found")
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

	s.Require().NoError(s.Vault.LoadIndex())

	_, err := s.Vault.LoadRecord(s.Record.Name, "")
	s.Assert().ErrorContains(err, fmt.Sprintf("name \"%s\" not found", s.Record.Name))
}
