package vault_test

import (
	"fmt"
	"pvault/data"
	"pvault/vault"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type OperationsTestSuite struct {
	suite.Suite
	Vault  vault.Vault
	Record vault.Record
}

func TestCreateCommandSuite(t *testing.T) {
	suite.Run(t, &OperationsTestSuite{})
}

func (s *OperationsTestSuite) SetupTest() {
	PATH := file.NewPath(s.T(), "vault")

	err := vault.InitializeNew(PATH)
	s.Require().NoError(err)

	s.Vault, err = vault.Open(PATH)
	s.Require().NoError(err)

	rand := rand.New(0)
	s.Record = vault.EmptyRecord(rand.ASCII(10))
	s.Record.Username = rand.ASCII(10)
	s.Record.Password = rand.ASCII(30)
	s.Record.Other = []interface{}{rand.ASCII(5), rand.ASCII(5), rand.ASCII(5)}
}

//=====================================

func (s *OperationsTestSuite) TestSaveRecordNameAlreadyExists() {
	//-- arrange
	s.Vault.Index[s.Record.Name] = uuid.Nil

	//-- act
	res := s.Vault.SaveRecord(s.Record)

	//-- assert
	s.Require().ErrorContains(res, fmt.Sprintf("name \"%s\" already exists", s.Record.Name))
}

func (s *OperationsTestSuite) TestSaveRecordInvalidVaultPath() {
	//-- arrange
	s.Vault.Path = "invalid"

	//-- act
	res := s.Vault.SaveRecord(s.Record)

	//-- assert
	s.Require().ErrorContains(res, "error saving record file")
}

func (s *OperationsTestSuite) TestSaveRecordNewNameValid() {
	//-- act
	res := s.Vault.SaveRecord(s.Record)

	//-- assert
	s.Require().NoError(res)
	s.Assert().Contains(s.Vault.Index, s.Record.Name)

	idx, err := vault.LoadIndex(s.Vault.IndexPath())
	s.Require().NoError(err)
	s.Assert().Contains(idx, s.Record.Name)

	r, err := data.LoadJSON[vault.Record](s.Vault.RecordPath(s.Record.ID))
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, r)
}

func (s *OperationsTestSuite) TestSaveRecordUpdateExistingValid() {
	//-- arrange
	s.Vault.Index[s.Record.Name] = s.Record.ID

	//-- act
	res := s.Vault.SaveRecord(s.Record)

	//-- assert
	s.Require().NoError(res)

	r, err := data.LoadJSON[vault.Record](s.Vault.RecordPath(s.Record.ID))
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, r)
}

func (s *OperationsTestSuite) TestLoadRecordInvalidVaultPath() {
	//-- arrange
	err := s.Vault.SaveRecord(s.Record)
	s.Require().NoError(err)

	s.Vault.Path = "invalid"

	//-- act
	_, res := s.Vault.LoadRecord(s.Record.ID)

	//-- assert
	s.Require().ErrorContains(res, "error loading record file")
}

func (s *OperationsTestSuite) TestLoadRecordValid() {
	//-- arrange
	err := s.Vault.SaveRecord(s.Record)
	s.Require().NoError(err)

	//-- act
	r, err := s.Vault.LoadRecord(s.Record.ID)

	//-- assert
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, r)
}
