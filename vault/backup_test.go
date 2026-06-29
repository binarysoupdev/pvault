package vault_test

import (
	"os"
	"path/filepath"
	"pvault/vault"
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
	s.Vault = vault.Vault{
		Path: file.NewPath(s.T(), ""),
	}
}

func (s *BackupTestSuite) TestBackupWhereVaultPathNotFoundReturnsError() {
	//-- arrange
	s.Vault.Path = "invalid"

	//-- act
	res := s.Vault.Backup("")

	//-- arrange
	s.Require().ErrorContains(res, "error reading vault directory")
}

func (s *BackupTestSuite) TestBackupWhereBackupDirAlreadyExistsReturnsError() {
	//-- arrange
	rand := rand.New(0)
	BACKUP := rand.ASCII(15)

	err := os.Mkdir(filepath.Join(s.Vault.Path, BACKUP), 0755)
	s.Require().NoError(err)

	//-- act
	res := s.Vault.Backup(BACKUP)

	//-- arrange
	s.Require().ErrorContains(res, "error creating backup directory")
}

func (s *BackupTestSuite) TestBackupValidCopiesOnlyFilesToBackupDir() {
	//-- arrange
	rand := rand.New(0)
	BACKUP := rand.ASCII(15)

	FILES := make([]string, 5)
	for i := range FILES {
		FILES[i] = rand.ASCII(10)
		os.WriteFile(filepath.Join(s.Vault.Path, FILES[i]), rand.Bytes(50), 0755)
	}

	DIR := rand.ASCII(10)
	err := os.Mkdir(filepath.Join(s.Vault.Path, DIR), 0755)
	s.Require().NoError(err)

	//-- act
	res := s.Vault.Backup(BACKUP)

	//-- arrange
	s.Require().NoError(res)

	backupDir := filepath.Join(s.Vault.Path, BACKUP)
	s.Require().DirExists(backupDir)

	for _, file := range FILES {
		bytes1, err := os.ReadFile(filepath.Join(s.Vault.Path, file))
		s.Require().NoError(err)

		bytes2, err := os.ReadFile(filepath.Join(backupDir, file))
		s.Require().NoError(err)

		s.Assert().Equal(bytes1, bytes2)
	}
	s.Assert().NoDirExists(filepath.Join(backupDir, DIR))
}
