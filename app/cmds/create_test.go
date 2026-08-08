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
	v2 "pvault/vault/record/record/v2"
	"testing"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type CreateTestSuite struct {
	test.CommandSuite[*cmds.CreateCommand]
	ConfigLoader json.Loader[config.Config]
	Config       config.Config

	Vault vault.Vault
}

func TestCreateCommandSuite(t *testing.T) {
	s := CreateTestSuite{
		ConfigLoader: json.NewLoader[config.Config](filepath.Join(t.TempDir(), "config.json")),
	}

	s.CommandSuite = test.NewCommandSuite(cmds.NewCreateCommand(s.ConfigLoader))
	suite.Run(t, &s)
}

func (s *CreateTestSuite) SetupTest() {
	s.Config = config.Config{
		Version:    config.VERSION,
		VaultPath:  filepath.Join(s.T().TempDir(), "vault"),
		OutputPath: s.T().TempDir(),
	}
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	var err error
	s.Vault, err = vault.InitializeNew(s.Config.VaultPath, "")
	s.Require().NoError(err)
}

//=====================================

func (s *CreateTestSuite) TestRunFailsWhenConfigNotFound() {
	//-- arrange
	s.Require().NoError(os.Remove(s.ConfigLoader.Path))

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("invalid config path")
}

func (s *CreateTestSuite) TestRunFailsWhenConfigVersionUnsupported() {
	//-- arrange
	s.Config.Version = config.VERSION + 1
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail(fmt.Sprintf("unsupported config version \"%d\"", s.Config.Version))
}

func (s *CreateTestSuite) TestRunFailsWhenNameIsEmpty() {
	//-- act
	s.RunCommand("-name", "")

	//-- assert
	s.RequireResultFail("\"name\" cannot be empty")
}

func (s *CreateTestSuite) TestRunFailsWithInvalidVaultPath() {
	//-- arrange
	const NAME = "name"

	s.Config.VaultPath = "invalid"
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	//-- act
	s.RunCommand("-name", NAME)

	//-- assert
	s.RequireResultFail("error opening vault")
}

func (s *CreateTestSuite) TestRunFailsWhenVaultOutOfDate() {
	//-- arrange
	const NAME = "name"

	s.Config.VaultPath = s.T().TempDir()
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	DATABASE := db_v1.Encoder{}

	META := meta.Metadata{
		DatabaseVersion: DATABASE.GetVersion(),
	}
	s.Require().NoError(meta.SaveMetadata(meta_v1.Encoder{}, s.Config.VaultPath, META))
	s.Require().NoError(database.SaveIndex(DATABASE, s.Config.VaultPath, index.IndexMap{}))

	//-- act
	s.RunCommand("-name", NAME)

	//-- assert
	s.RequireResultFail(fmt.Sprintf("vault (@v%d) out-of-date", DATABASE.GetVersion()))
}

func (s *CreateTestSuite) TestRunFailsWhenConfigOutputPathInvalid() {
	//-- arrange
	const NAME = "name"

	s.Config.OutputPath = "invalid"
	s.Require().NoError(json.MarshalFile(s.Config, s.ConfigLoader.Path))

	//-- act
	s.RunCommand("-name", NAME)

	//-- assert
	s.RequireResultFail("error validating config \"output_path\"")
}

func (s *CreateTestSuite) TestRunFailsWhenNameAlreadyExists() {
	//-- arrange
	const NAME = "name"
	const PASSWORD = "Password123!"

	v, err := vault.Open(s.Config.VaultPath)
	s.Require().NoError(err)

	s.Require().NoError(v.SaveRecord(v2.NewEmptyRecord(NAME), PASSWORD))

	//-- act
	s.RunCommand("-name", NAME)

	//-- assert
	s.RequireResultFail("error validating record")
}

func (s *CreateTestSuite) TestRunFailsWithIncorrectVerifyPassword() {
	//-- arrange
	const NAME = "name"
	const PASSWORD = "Password123!"

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

func (s *CreateTestSuite) TestRunPassesAndCreatesNewRecord() {
	//-- arrange
	const NAME = "name"
	const PASSWORD = "Password123!"

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
	s.Require().Contains(line, "[+] Saved Record: ")

	ID, err := uuid.Parse(line[len(line)-36:])
	s.Require().NoError(err)

	OUTPUT_FILE := filepath.Join(s.Config.OutputPath, ID.String()+".json")
	s.Assert().Contains(io.ReadLine(), "[+] "+OUTPUT_FILE)

	s.Require().NoError(s.Vault.LoadIndex())

	r1, err := json.UnmarshalFile[v2.Record](OUTPUT_FILE)
	s.Require().NoError(err)

	r2, err := s.Vault.LoadRecord(NAME, PASSWORD)
	s.Require().NoError(err)

	s.Assert().Equal(r1, r2)
}
