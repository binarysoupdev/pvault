package vault_test

import (
	"pvault/vault"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type OperationsTestSuite struct {
	suite.Suite
	Vault  vault.Vault
	Record vault.Record
}

func (s *OperationsTestSuite) SetupTest() {
	PATH := file.NewPath(s.T(), "vault")

	err := vault.InitializeNew(PATH)
	s.Require().NoError(err)

	s.Vault, err = vault.Open(PATH)
	s.Require().NoError(err)

	rand := rand.New(0)
	s.Record = vault.NewRecord(rand.ASCII(10))
}

func (s *OperationsTestSuite) TestSaveRecordInvalidVaultPath(t *testing.T) {
	//-- arrange
	s.Vault.Path = "invalid"

	//-- act
	res := s.Vault.SaveRecord(s.Record)

	//-- assert
	s.Require().ErrorContains(res, "error saving index file")
}

func (s *OperationsTestSuite) TestSaveRecordValid(t *testing.T) {
	//-- act
	res := s.Vault.SaveRecord(s.Record)

	//-- assert
	s.Require().NoError(res)
	s.Assert().Contains(s.Vault.Index, s.Record.Name)

	//TODO: load and verify index file
	//TODO: load and verify vault record file
}
