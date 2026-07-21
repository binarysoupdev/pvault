package local_test

import (
	"fmt"
	"os"
	"pvault/app/vault/index"
	"pvault/app/vault/local"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type BackupTestSuite struct {
	suite.Suite
	Vault local.Vault
}

func TestBackupTestSuite(t *testing.T) {
	suite.Run(t, &BackupTestSuite{})
}

func (s *BackupTestSuite) SetupTest() {
	s.Vault = local.Vault{
		Index: &index.Mock{
			Path: file.NewPath(s.T(), ""),
		},
		Map: index.IndexMap{},
	}
}

//===================================

func (s *BackupTestSuite) TestBackupReturnsErrorWhenPathNotFound() {
	//-- act
	res := s.Vault.Backup("invalid")

	//-- arrange
	s.Require().ErrorContains(res, "error reading backup directory")
}

func (s *BackupTestSuite) TestBackupReturnsErrorWhenPathIsNotADir() {
	//-- arrange
	PATH := file.CreateEmpty(s.T(), "backups.txt")

	//-- act
	res := s.Vault.Backup(PATH)

	//-- arrange
	s.Require().ErrorContains(res, fmt.Sprintf("\"%s\" is not a directory", PATH))
}

func (s *BackupTestSuite) TestBackupReturnsNoErrorAndBacksUpIndexFileAndRecords() {
	//-- arrange
	BACKUP := file.NewPath(s.T(), "")
	IDS := uuid.UUIDs{uuid.New(), uuid.New()}
	DATA := [][]byte{{0, 0, 0}, {1, 1, 1}, {2, 2, 2}}

	s.Require().NoError(os.WriteFile(s.Vault.Index.Filepath(), DATA[0], 0666))
	s.Require().NoError(os.WriteFile(s.Vault.Index.RecordPath(IDS[0]), DATA[1], 0666))
	s.Require().NoError(os.WriteFile(s.Vault.Index.RecordPath(IDS[1]), DATA[2], 0666))

	s.Vault.Map = index.IndexMap{
		"foo1": IDS[0],
		"foo2": IDS[1],
	}

	//-- act
	res := s.Vault.Backup(BACKUP)

	//-- arrange
	s.Require().NoError(res)
	idx := index.Mock{
		Path: BACKUP,
	}

	data0, err := os.ReadFile(idx.Filepath())
	s.Require().NoError(err)
	s.Assert().Equal(DATA[0], data0)

	data1, err := os.ReadFile(idx.RecordPath(IDS[0]))
	s.Require().NoError(err)
	s.Assert().Equal(DATA[1], data1)

	data2, err := os.ReadFile(idx.RecordPath(IDS[1]))
	s.Require().NoError(err)
	s.Assert().Equal(DATA[2], data2)
}
