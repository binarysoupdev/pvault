package v2_test

import (
	"os"
	record "pvault/vault/record/version2"
	v2 "pvault/vault/record/version2"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type RecordSuite struct {
	suite.Suite
	Record v2.Record
}

func TestRecordSuite(t *testing.T) {
	suite.Run(t, &RecordSuite{})
}

func (s *RecordSuite) SetupTest() {
	s.Record = record.Record{
		ID:       uuid.New(),
		Name:     "name",
		Username: "username",
		Password: "password",
		//Other:    map[string]any{},
	}
}

//=====================================

func (s *RecordSuite) TestValidateReturnsErrorWhenInvalid() {
	//-- arrange
	s.Record.ID = uuid.Nil
	s.Record.Name = ""

	//-- act
	res := s.Record.Validate()

	//-- assert
	s.Assert().ErrorContains(res, "id cannot be nil (all zeroes)")
	s.Assert().ErrorContains(res, "name cannot be empty")
}

func (s *RecordSuite) TestValidateReturnsNoErrorWhenValid() {
	//-- act
	res := s.Record.Validate()

	//-- assert
	s.Require().NoError(res)
}

func (s *RecordSuite) TestSaveFileReturnsNoErrorAndEncodesRecordToFile() {
	//-- arrange
	PATH := file.NewPath(s.T(), "record.json")
	const PASSWORD = "Password123!"

	//-- act
	err := s.Record.SaveFile(PATH, PASSWORD)

	//-- assert
	s.Require().NoError(err)

	bytes, err := os.ReadFile(PATH)
	s.Require().NoError(err)

	res, err := v2.Unmarshal(bytes, PASSWORD)
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, res)
}
