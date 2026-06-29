package vault_test

import (
	"fmt"
	"pvault/vault"
	"pvault/vault/record"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type DatabaseTestSuite struct {
	suite.Suite
	Vault    vault.Vault
	Record   record.Record
	Password string
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, &DatabaseTestSuite{})
}

func (s *DatabaseTestSuite) SetupTest() {
	PATH := file.NewPath(s.T(), "vault")

	_, err := vault.InitializeNew(PATH)
	s.Require().NoError(err)

	s.Vault, err = vault.Open(PATH)
	s.Require().NoError(err)

	rand := rand.New(0)
	s.Record = record.NewFromName(rand.ASCII(10))
	s.Record.Username = rand.ASCII(10)
	s.Record.Password = rand.ASCII(30)
	s.Record.Other = []interface{}{rand.ASCII(5), rand.ASCII(5), rand.ASCII(5)}

	s.Password = rand.ASCII(30)
}

//=====================================

func (s *DatabaseTestSuite) TestSaveRecordInvalidVaultPath() {
	//-- arrange
	s.Vault.Path = "invalid"

	//-- act
	res := s.Vault.SaveRecord(s.Record, s.Password)

	//-- assert
	s.Require().ErrorContains(res, "error writing record file")
}

func (s *DatabaseTestSuite) TestSaveRecordNewIdNewNameValid() {
	//-- act
	res := s.Vault.SaveRecord(s.Record, s.Password)

	//-- assert
	s.Require().NoError(res)
	s.Assert().Contains(s.Vault.Index, s.Record.Name)

	idx, err := s.Vault.Database.LoadIndex()
	s.Require().NoError(err)
	s.Assert().Contains(idx, s.Record.Name)

	other, err := s.Vault.LoadRecord(s.Record.Name, s.Password)
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, other)
}

func (s *DatabaseTestSuite) TestSaveRecordExistingIDNewNameValid() {
	//-- arrange
	s.Vault.Index[s.Record.Name] = s.Record.ID
	s.Record.Name += "x"

	//-- act
	res := s.Vault.SaveRecord(s.Record, s.Password)

	//-- assert
	s.Require().NoError(res)
	s.Assert().Contains(s.Vault.Index, s.Record.Name)
	s.Assert().Len(s.Vault.Index, 1)

	idx, err := s.Vault.Database.LoadIndex()
	s.Require().NoError(err)
	s.Assert().Contains(idx, s.Record.Name)
	s.Assert().Len(s.Vault.Index, 1)

	other, err := s.Vault.LoadRecord(s.Record.Name, s.Password)
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, other)
}

func (s *DatabaseTestSuite) TestSaveRecordNewIDExistingNameInvalid() {
	//-- arrange
	s.Vault.Index[s.Record.Name] = uuid.Nil

	//-- act
	res := s.Vault.SaveRecord(s.Record, s.Password)

	//-- assert
	s.Require().ErrorContains(res, "error validating record")
}

func (s *DatabaseTestSuite) TestSaveRecordExistingIDExistingNameValid() {
	//-- arrange
	s.Vault.Index[s.Record.Name] = s.Record.ID

	//-- act
	res := s.Vault.SaveRecord(s.Record, s.Password)

	//-- assert
	s.Require().NoError(res)
	s.Assert().Len(s.Vault.Index, 1)

	idx, err := s.Vault.Database.LoadIndex()
	s.Require().NoError(err)
	s.Assert().Len(idx, 1)

	other, err := s.Vault.LoadRecord(s.Record.Name, s.Password)
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, other)
}

func (s *DatabaseTestSuite) TestLoadRecordInvalidVaultPath() {
	//-- arrange
	err := s.Vault.SaveRecord(s.Record, s.Password)
	s.Require().NoError(err)

	s.Vault.Path = "invalid"

	//-- act
	_, res := s.Vault.LoadRecord(s.Record.Name, s.Password)

	//-- assert
	s.Require().ErrorContains(res, "error reading record file")
}

func (s *DatabaseTestSuite) TestLoadRecordInvalidName() {
	//-- arrange
	err := s.Vault.SaveRecord(s.Record, s.Password)
	s.Require().NoError(err)

	NAME := s.Record.Name + "x"

	//-- act
	_, res := s.Vault.LoadRecord(NAME, s.Password)

	//-- assert
	s.Require().ErrorContains(res, fmt.Sprintf("name \"%s\" not found", NAME))
}

func (s *DatabaseTestSuite) TestLoadRecordIncorrectPassword() {
	//-- arrange
	err := s.Vault.SaveRecord(s.Record, s.Password)
	s.Require().NoError(err)

	//-- act
	_, res := s.Vault.LoadRecord(s.Record.Name, s.Password+"x")

	//-- assert
	s.Require().ErrorContains(res, "error decrypting ciphertext")
}

func (s *DatabaseTestSuite) TestLoadRecordValid() {
	//-- arrange
	err := s.Vault.SaveRecord(s.Record, s.Password)
	s.Require().NoError(err)

	//-- act
	r, err := s.Vault.LoadRecord(s.Record.Name, s.Password)

	//-- assert
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, r)
}

func (s *DatabaseTestSuite) TestDeleteRecordInvalidName() {
	//-- act
	_, res := s.Vault.DeleteRecord(s.Record.Name)

	//-- assert
	s.Require().ErrorContains(res, fmt.Sprintf("name \"%s\" not found", s.Record.Name))
}

func (s *DatabaseTestSuite) TestDeleteRecordInvalidVaultPath() {
	//-- arrange
	err := s.Vault.SaveRecord(s.Record, s.Password)
	s.Require().NoError(err)

	s.Vault.Path = "invalid"

	//-- act
	_, res := s.Vault.DeleteRecord(s.Record.Name)

	//-- assert
	s.Require().ErrorContains(res, "error deleting record file")
}

func (s *DatabaseTestSuite) TestDeleteRecordValid() {
	//-- arrange
	err := s.Vault.SaveRecord(s.Record, s.Password)
	s.Require().NoError(err)

	//-- act
	id, err := s.Vault.DeleteRecord(s.Record.Name)

	//-- assert
	s.Require().NoError(err)
	s.Assert().Equal(s.Record.ID, id)

	idx, err := s.Vault.Database.LoadIndex()
	s.Require().NoError(err)
	s.Assert().Len(idx, 0)

	s.Assert().NoFileExists(s.Database.RecordPath(s.Record.ID))
}
