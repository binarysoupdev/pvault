package v1_test

import (
	"os"
	record "pvault/vault/record/version1"
	v1 "pvault/vault/record/version1"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type RecordSuite struct {
	suite.Suite
	Record v1.Record
}

func TestRecordSuite(t *testing.T) {
	suite.Run(t, &RecordSuite{})
}

func (s *RecordSuite) SetupTest() {
	s.Record = record.Record{
		Password:      "password",
		Username:      "username",
		URL:           "url",
		RecoveryCodes: []string{"code1", "code2"},
		ID:            uuid.New(),
		Name:          "name",
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

	res, err := v1.Unmarshal(bytes, PASSWORD, s.Record.ID)
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, res)
}

func (s *RecordSuite) TestUpgradeReturnsUpgradedRecord() {
	//-- act
	res := s.Record.Upgrade()

	//-- assert
	s.Assert().Equal(s.Record.ID, res.ID)
	s.Assert().Equal(s.Record.Name, res.Name)
	s.Assert().Equal(s.Record.Password, res.Password)
	s.Assert().Equal(s.Record.Username, res.Username)
	s.Assert().Equal(s.Record.URL, res.Other["url"])
	s.Assert().Equal(s.Record.RecoveryCodes, res.Other["recovery_codes"])
}
