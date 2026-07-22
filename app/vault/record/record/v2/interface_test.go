package v2_test

import (
	"os"
	record "pvault/app/vault/record/version2"
	v2 "pvault/app/vault/record/version2"
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
	}
}

//=====================================

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
