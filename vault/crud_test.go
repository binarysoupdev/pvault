package vault_test

import (
	"fmt"
	"pvault/data"
	"pvault/vault"
	"pvault/vault/index"
	"pvault/vault/record"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type CRUDTestSuite struct {
	suite.Suite
	Vault  vault.Vault
	Record record.Record
}

func TestCreateCommandSuite(t *testing.T) {
	suite.Run(t, &CRUDTestSuite{})
}

func (s *CRUDTestSuite) SetupTest() {
	PATH := file.NewPath(s.T(), "vault")

	err := vault.InitializeNew(PATH)
	s.Require().NoError(err)

	s.Vault, err = vault.Open(PATH)
	s.Require().NoError(err)

	rand := rand.New(0)
	s.Record = record.NewFromName(rand.ASCII(10))
	s.Record.Username = rand.ASCII(10)
	s.Record.Password = rand.ASCII(30)
	s.Record.Other = []interface{}{rand.ASCII(5), rand.ASCII(5), rand.ASCII(5)}
}

//=====================================

func (s *CRUDTestSuite) TestSaveRecordInvalidVaultPath() {
	//-- arrange
	s.Vault.Path = "invalid"

	//-- act
	res := s.Vault.SaveRecord(s.Record)

	//-- assert
	s.Require().ErrorContains(res, "error saving record file")
}

func (s *CRUDTestSuite) TestSaveRecordNewIdNewNameValid() {
	//-- act
	res := s.Vault.SaveRecord(s.Record)

	//-- assert
	s.Require().NoError(res)
	s.Assert().Contains(s.Vault.Index, s.Record.Name)

	idx, err := index.LoadIndex(s.Vault.IndexPath())
	s.Require().NoError(err)
	s.Assert().Contains(idx, s.Record.Name)

	r, err := data.LoadJSON[record.Record](s.Vault.RecordPath(s.Record.ID))
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, r)
}

func (s *CRUDTestSuite) TestSaveRecordExistingIDNewNameValid() {
	//-- arrange
	s.Vault.Index[s.Record.Name] = s.Record.ID
	s.Record.Name += "x"

	//-- act
	res := s.Vault.SaveRecord(s.Record)

	//-- assert
	s.Require().NoError(res)
	s.Assert().Contains(s.Vault.Index, s.Record.Name)
	s.Assert().Len(s.Vault.Index, 1)

	idx, err := index.LoadIndex(s.Vault.IndexPath())
	s.Require().NoError(err)
	s.Assert().Contains(idx, s.Record.Name)
	s.Assert().Len(s.Vault.Index, 1)

	r, err := data.LoadJSON[record.Record](s.Vault.RecordPath(s.Record.ID))
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, r)
}

func (s *CRUDTestSuite) TestSaveRecordNewIDExistingNameInvalid() {
	//-- arrange
	s.Vault.Index[s.Record.Name] = uuid.Nil

	//-- act
	res := s.Vault.SaveRecord(s.Record)

	//-- assert
	s.Require().ErrorContains(res, fmt.Sprintf("name \"%s\" already exists", s.Record.Name))
}

func (s *CRUDTestSuite) TestSaveRecordExistingIDExistingNameValid() {
	//-- arrange
	s.Vault.Index[s.Record.Name] = s.Record.ID

	//-- act
	res := s.Vault.SaveRecord(s.Record)

	//-- assert
	s.Require().NoError(res)
	s.Assert().Len(s.Vault.Index, 1)

	idx, err := index.LoadIndex(s.Vault.IndexPath())
	s.Require().NoError(err)
	s.Assert().Len(idx, 1)

	r, err := data.LoadJSON[record.Record](s.Vault.RecordPath(s.Record.ID))
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, r)
}

func (s *CRUDTestSuite) TestLoadRecordInvalidVaultPath() {
	//-- arrange
	err := s.Vault.SaveRecord(s.Record)
	s.Require().NoError(err)

	s.Vault.Path = "invalid"

	//-- act
	_, res := s.Vault.LoadRecord(s.Record.Name)

	//-- assert
	s.Require().ErrorContains(res, "error loading record file")
}

func (s *CRUDTestSuite) TestLoadRecordInvalidName() {
	//-- arrange
	err := s.Vault.SaveRecord(s.Record)
	s.Require().NoError(err)

	NAME := s.Record.Name + "x"

	//-- act
	_, res := s.Vault.LoadRecord(NAME)

	//-- assert
	s.Require().ErrorContains(res, fmt.Sprintf("name \"%s\" not found", NAME))
}

func (s *CRUDTestSuite) TestLoadRecordValid() {
	//-- arrange
	err := s.Vault.SaveRecord(s.Record)
	s.Require().NoError(err)

	//-- act
	r, err := s.Vault.LoadRecord(s.Record.Name)

	//-- assert
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, r)
}

func (s *CRUDTestSuite) TestDeleteRecordInvalidName() {
	//-- act
	_, res := s.Vault.DeleteRecord(s.Record.Name)

	//-- assert
	s.Require().ErrorContains(res, fmt.Sprintf("name \"%s\" not found", s.Record.Name))
}

func (s *CRUDTestSuite) TestDeleteRecordInvalidVaultPath() {
	//-- arrange
	err := s.Vault.SaveRecord(s.Record)
	s.Require().NoError(err)

	s.Vault.Path = "invalid"

	//-- act
	_, res := s.Vault.DeleteRecord(s.Record.Name)

	//-- assert
	s.Require().ErrorContains(res, "error deleting record file")
}

func (s *CRUDTestSuite) TestDeleteRecordValid() {
	//-- arrange
	err := s.Vault.SaveRecord(s.Record)
	s.Require().NoError(err)

	//-- act
	id, err := s.Vault.DeleteRecord(s.Record.Name)

	//-- assert
	s.Require().NoError(err)
	s.Assert().Equal(s.Record.ID, id)

	idx, err := index.LoadIndex(s.Vault.IndexPath())
	s.Require().NoError(err)
	s.Assert().Len(idx, 0)

	s.Assert().NoFileExists(s.Vault.RecordPath(s.Record.ID))
}
