package vault_test

import (
	"fmt"
	"pvault/vault"
	v2 "pvault/vault/record/version/v2"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type BackupTestSuite struct {
	suite.Suite
	Vault vault.Vault
}

func TestBackupTestSuite(t *testing.T) {
	suite.Run(t, &BackupTestSuite{})
}

func (s *BackupTestSuite) SetupTest() {
	var err error
	s.Vault, err = vault.InitializeNew(file.NewPath(s.T(), "vault"))
	s.Require().NoError(err)
}

func (s *BackupTestSuite) TestBackupWherePathNotFoundReturnsError() {
	//-- act
	res := s.Vault.Backup("invalid")

	//-- arrange
	s.Require().ErrorContains(res, "error reading backup directory")
}

func (s *BackupTestSuite) TestBackupWherePathIsNotADirReturnsError() {
	//-- arrange
	PATH := file.CreateEmpty(s.T(), "backups.txt")

	//-- act
	res := s.Vault.Backup(PATH)

	//-- arrange
	s.Require().ErrorContains(res, fmt.Sprintf("\"%s\" is not a directory", PATH))
}

func (s *BackupTestSuite) TestBackupValidBacksUpIndexFileAndRecord() {
	//-- arrange
	BACKUP := file.NewPath(s.T(), "")

	rand := rand.New(0)
	const NUM_RECORDS = 5
	for range NUM_RECORDS {
		err := s.Vault.SaveRecord(v2.NewFromName(rand.ASCII(10)), rand.ASCII(30))
		s.Require().NoError(err)
	}

	//-- act
	res := s.Vault.Backup(BACKUP)

	//-- arrange
	s.Require().NoError(res)

	v, err := vault.Open(BACKUP)
	s.Require().NoError(err)

	s.Assert().Equal(s.Vault.Index, v.Index)
}
