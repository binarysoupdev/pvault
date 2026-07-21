package vault_test

import (
	"fmt"
	"pvault/app/vault/index"
	"pvault/app/vault/local"
	"pvault/app/vault/record"
	record_v1 "pvault/app/vault/record/version1"
	record_v2 "pvault/app/vault/record/version2"
	"testing"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type RecordTestSuite struct {
	suite.Suite
	Index  index.Mock
	Record *record.Mock
	Vault  local.Vault
}

func TestRecordTestSuite(t *testing.T) {
	suite.Run(t, &RecordTestSuite{})
}

func (s *RecordTestSuite) SetupTest() {
	s.Index = index.Mock{
		Path: file.NewPath(s.T(), ""),
	}

	s.Record = &record.Mock{
		ID:   uuid.New(),
		Name: "name",
	}

	s.Vault = local.Vault{
		Index: &s.Index,
		Map:   index.IndexMap{},
	}
}

//===================================

func (s *RecordTestSuite) TestValidateRecordReturnsErrorWhenRecordInvalid() {
	//-- act
	res := s.Vault.ValidateRecord(&record.Mock{})

	//-- assert
	s.Require().ErrorContains(res, "record invalid")
}

func (s *RecordTestSuite) TestValidateRecordReturnsErrorWhenNameAlreadyExistsForAnotherRecord() {
	//-- arrange
	s.Vault.Map = index.IndexMap{
		s.Record.Name: uuid.New(),
	}

	//-- act
	res := s.Vault.ValidateRecord(s.Record)

	//-- assert
	s.Require().ErrorContains(res, fmt.Sprintf("name \"%s\" already exists", s.Record.Name))
}

func (s *RecordTestSuite) TestValidateRecordReturnsNoErrorWhenNameExistsForSameRecord() {
	//-- arrange
	s.Vault.Map = index.IndexMap{
		s.Record.Name: s.Record.ID,
	}

	//-- act
	res := s.Vault.ValidateRecord(s.Record)

	//-- assert
	s.Require().NoError(res)
}

func (s *RecordTestSuite) TestSaveRecordReturnsErrorWhenRecordInvalid() {
	//-- act
	res := s.Vault.SaveRecord(&record.Mock{}, "")

	//-- assert
	s.Require().ErrorContains(res, "error validating record")
}

func (s *RecordTestSuite) TestSaveRecordReturnsErrorWhenRecordSaveFileReturnsError() {
	//-- arrange
	s.Record.SaveFileError = errors.New("")

	//-- act
	res := s.Vault.SaveRecord(s.Record, "")

	//-- assert
	s.Require().ErrorContains(res, "error saving record")
}

func (s *RecordTestSuite) TestSaveRecordReturnsErrorWhenIndexSaveMapReturnsError() {
	//-- arrange
	s.Index.SaveMapError = errors.New("")

	//-- act
	res := s.Vault.SaveRecord(s.Record, "")

	//-- assert
	s.Require().ErrorContains(res, "error saving index map")
}

func (s *RecordTestSuite) TestSaveRecordReturnsNoErrorWithNewIdAndNewName() {
	//-- act
	res := s.Vault.SaveRecord(s.Record, "")

	//-- assert
	s.Require().NoError(res)
	s.Assert().Contains(s.Index.Map, s.Record.Name)
}

func (s *RecordTestSuite) TestSaveRecordReturnsNoErrorWithExistingIDAndNewName() {
	//-- arrange
	s.Vault.Map[s.Record.Name] = s.Record.ID
	s.Record.Name += "x"

	//-- act
	res := s.Vault.SaveRecord(s.Record, "")

	//-- assert
	s.Require().NoError(res)
	s.Assert().Contains(s.Index.Map, s.Record.Name)
	s.Assert().Len(s.Index.Map, 1)
}

func (s *RecordTestSuite) TestSaveRecordReturnsErrorWithNewIDAndExistingName() {
	//-- arrange
	s.Vault.Map[s.Record.Name] = uuid.New()

	//-- act
	res := s.Vault.SaveRecord(s.Record, "")

	//-- assert
	s.Require().ErrorContains(res, "error validating record")
}

func (s *RecordTestSuite) TestSaveRecordReturnsNoErrorWithExistingIDAndExistingName() {
	//-- arrange
	s.Vault.Map[s.Record.Name] = s.Record.ID

	//-- act
	res := s.Vault.SaveRecord(s.Record, "")

	//-- assert
	s.Require().NoError(res)
	s.Assert().Contains(s.Index.Map, s.Record.Name)
	s.Assert().Len(s.Index.Map, 1)
}

func (s *RecordTestSuite) TestLoadRecordReturnsErrorWhenNameNotFound() {
	//-- act
	_, res := s.Vault.LoadRecord(s.Record.Name, "")

	//-- assert
	s.Require().ErrorContains(res, fmt.Sprintf("name \"%s\" not found", s.Record.Name))
}

func (s *RecordTestSuite) TestLoadRecordReturnsErrorWhenRecordNotFound() {
	//-- arrange
	s.Vault.Map[s.Record.Name] = s.Record.ID

	//-- act
	_, res := s.Vault.LoadRecord(s.Record.Name, "")

	//-- assert
	s.Require().ErrorContains(res, "error loading record")
}

func (s *RecordTestSuite) TestLoadRecordReturnsV1RecordAndNoError() {
	//-- arrange
	RECORD := record_v1.Record{
		ID:   uuid.New(),
		Name: "name",
	}
	const PASSWORD = "Password123!"

	err := s.Vault.SaveRecord(RECORD, PASSWORD)
	s.Require().NoError(err)

	//-- act
	res, err := s.Vault.LoadRecord(RECORD.Name, PASSWORD)

	//-- assert
	s.Require().NoError(err)
	s.Assert().Equal(RECORD, res)
}

func (s *RecordTestSuite) TestLoadRecordReturnsV2RecordAndNoError() {
	//-- arrange
	RECORD := record_v2.Record{
		ID:   uuid.New(),
		Name: "name",
	}
	const PASSWORD = "Password123!"

	err := s.Vault.SaveRecord(RECORD, PASSWORD)
	s.Require().NoError(err)

	//-- act
	res, err := s.Vault.LoadRecord(RECORD.Name, PASSWORD)

	//-- assert
	s.Require().NoError(err)
	s.Assert().Equal(RECORD, res)
}

func (s *RecordTestSuite) TestDeleteRecordReturnsErrorWhenNameNotFound() {
	//-- act
	_, res := s.Vault.DeleteRecord(s.Record.Name)

	//-- assert
	s.Require().ErrorContains(res, fmt.Sprintf("name \"%s\" not found", s.Record.Name))
}

func (s *RecordTestSuite) TestDeleteRecordReturnsErrorWhenRecordNotFound() {
	//-- arrange
	s.Vault.Map[s.Record.Name] = s.Record.ID

	//-- act
	_, res := s.Vault.DeleteRecord(s.Record.Name)

	//-- assert
	s.Require().ErrorContains(res, "error deleting record")
}

func (s *RecordTestSuite) TestDeleteRecordReturnsErrorWhenIndexSaveMapReturnsError() {
	//-- arrange
	const PASSWORD = "Password123!"
	err := s.Vault.SaveRecord(s.Record, PASSWORD)
	s.Require().NoError(err)

	s.Index.SaveMapError = errors.New("")

	//-- act
	_, res := s.Vault.DeleteRecord(s.Record.Name)

	//-- assert
	s.Require().ErrorContains(res, "error saving index map")
}

func (s *RecordTestSuite) TestDeleteRecordReturnsIDAndNoErrorAndDeletesRecord() {
	//-- arrange
	const PASSWORD = "Password123!"
	err := s.Vault.SaveRecord(s.Record, PASSWORD)
	s.Require().NoError(err)

	//-- act
	id, err := s.Vault.DeleteRecord(s.Record.Name)

	//-- assert
	s.Require().NoError(err)

	s.Assert().Equal(s.Record.ID, id)
	s.Assert().Empty(s.Index.Map)
	s.Assert().NoFileExists(s.Index.RecordPath(id))
}
