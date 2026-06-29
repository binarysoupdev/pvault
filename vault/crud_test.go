package vault_test

import (
	"errors"
	"fmt"
	"pvault/vault"
	"pvault/vault/data"
	"pvault/vault/index"
	"pvault/vault/record"
	"testing"

	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type CRUDTestSuite struct {
	suite.Suite
	DatabaseMock data.DatabaseMock

	Vault    vault.Vault
	Record   record.Record
	Password string
}

func TestCRUDTestSuite(t *testing.T) {
	suite.Run(t, &CRUDTestSuite{})
}

func (s *CRUDTestSuite) SetupTest() {
	s.DatabaseMock = data.DatabaseMock{}

	s.Vault = vault.Vault{
		Index:    index.IndexMap{},
		Database: &s.DatabaseMock,
	}

	rand := rand.New(0)
	s.Record = record.NewFromName(rand.ASCII(10))
	s.Record.Username = rand.ASCII(10)
	s.Record.Password = rand.ASCII(30)
	s.Record.Other = []interface{}{rand.ASCII(5), rand.ASCII(5), rand.ASCII(5)}
}

//=====================================

func (s *CRUDTestSuite) TestSaveRecordWithInvalidRecordReturnsError() {
	//-- act
	res := s.Vault.SaveRecord(record.Record{}, "")

	//-- assert
	s.Require().ErrorContains(res, "error validating record")
}

func (s *CRUDTestSuite) TestSaveRecordWhereDatabaseSaveRecordFailsReturnsError() {
	//-- arrange
	s.DatabaseMock.SaveRecordError = errors.New("")

	//-- act
	res := s.Vault.SaveRecord(s.Record, "")

	//-- assert
	s.Require().ErrorContains(res, "error saving record to database")
}

func (s *CRUDTestSuite) TestSaveRecordWhereDatabaseSaveIndexFailsReturnsError() {
	//-- arrange
	s.DatabaseMock.SaveIndexError = errors.New("")

	//-- act
	res := s.Vault.SaveRecord(s.Record, "")

	//-- assert
	s.Require().ErrorContains(res, "error saving index to database")
}

func (s *CRUDTestSuite) TestSaveRecordWithNewIdAndNewNameIsValid() {
	//-- act
	res := s.Vault.SaveRecord(s.Record, "")

	//-- assert
	s.Require().NoError(res)
	s.Assert().Equal(s.Record, s.DatabaseMock.Record)

	s.Assert().Contains(s.Vault.Index, s.Record.Name)
	s.Assert().Contains(s.DatabaseMock.Index, s.Record.Name)

}

func (s *CRUDTestSuite) TestSaveRecordWithExistingIDAndNewNameIsValid() {
	//-- arrange
	s.Vault.Index[s.Record.Name] = s.Record.ID
	s.Record.Name += "x"

	//-- act
	res := s.Vault.SaveRecord(s.Record, "")

	//-- assert
	s.Require().NoError(res)
	s.Assert().Equal(s.Record, s.DatabaseMock.Record)

	s.Assert().Contains(s.Vault.Index, s.Record.Name)
	s.Assert().Len(s.Vault.Index, 1)

	s.Assert().Contains(s.DatabaseMock.Index, s.Record.Name)
	s.Assert().Len(s.DatabaseMock.Index, 1)
}

func (s *CRUDTestSuite) TestSaveRecordWithNewIDAndExistingNameIsInvalid() {
	//-- arrange
	s.Vault.Index[s.Record.Name] = uuid.Nil

	//-- act
	res := s.Vault.SaveRecord(s.Record, "")

	//-- assert
	s.Require().ErrorContains(res, "error validating record")
}

func (s *CRUDTestSuite) TestSaveRecordWithExistingIDAndExistingNameIsValid() {
	//-- arrange
	s.Vault.Index[s.Record.Name] = s.Record.ID

	//-- act
	res := s.Vault.SaveRecord(s.Record, "")

	//-- assert
	s.Require().NoError(res)
	s.Assert().Equal(s.Record, s.DatabaseMock.Record)

	s.Assert().Contains(s.Vault.Index, s.Record.Name)
	s.Assert().Len(s.Vault.Index, 1)

	s.Assert().Contains(s.DatabaseMock.Index, s.Record.Name)
	s.Assert().Len(s.DatabaseMock.Index, 1)
}

func (s *CRUDTestSuite) TestLoadRecordWhereNameNotFoundReturnsError() {
	//-- act
	_, res := s.Vault.LoadRecord(s.Record.Name, "")

	//-- assert
	s.Require().ErrorContains(res, fmt.Sprintf("name \"%s\" not found", s.Record.Name))
}

func (s *CRUDTestSuite) TestLoadRecordWhereDatabaseLoadRecordFailsReturnsError() {
	//-- arrange
	err := s.Vault.SaveRecord(s.Record, "")
	s.Require().NoError(err)

	s.DatabaseMock.LoadRecordError = errors.New("")

	//-- act
	_, res := s.Vault.LoadRecord(s.Record.Name, "")

	//-- assert
	s.Require().ErrorContains(res, "error loading record from database")
}

func (s *CRUDTestSuite) TestLoadRecordValidReturnsRecord() {
	//-- arrange
	err := s.Vault.SaveRecord(s.Record, "")
	s.Require().NoError(err)

	//-- act
	r, err := s.Vault.LoadRecord(s.Record.Name, "")

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

func (s *CRUDTestSuite) TestDeleteRecordWhereDatabaseDeleteRecordFailsReturnsError() {
	//-- arrange
	err := s.Vault.SaveRecord(s.Record, "")
	s.Require().NoError(err)

	s.DatabaseMock.DeleteRecordError = errors.New("")

	//-- act
	_, res := s.Vault.DeleteRecord(s.Record.Name)

	//-- assert
	s.Require().ErrorContains(res, "error deleting record from database")
}

func (s *CRUDTestSuite) TestDeleteRecordWhereDatabaseSaveIndexFailsReturnsError() {
	//-- arrange
	err := s.Vault.SaveRecord(s.Record, "")
	s.Require().NoError(err)

	s.DatabaseMock.SaveIndexError = errors.New("")

	//-- act
	_, res := s.Vault.DeleteRecord(s.Record.Name)

	//-- assert
	s.Require().ErrorContains(res, "error saving index to database")
}

func (s *CRUDTestSuite) TestDeleteRecordValidDeletesRecord() {
	//-- arrange
	err := s.Vault.SaveRecord(s.Record, "")
	s.Require().NoError(err)

	//-- act
	id, err := s.Vault.DeleteRecord(s.Record.Name)

	//-- assert
	s.Require().NoError(err)
	s.Assert().Equal(s.Record.ID, id)

	s.Assert().Empty(s.Vault.Index)
	s.Assert().Empty(s.DatabaseMock.Index)
	s.Assert().NotEqual(s.Record, s.DatabaseMock.Record)
}
