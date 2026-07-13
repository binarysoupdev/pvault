package version2_test

import (
	"encoding/binary"
	"fmt"
	"os"
	"pvault/vault/database/version2"
	v1 "pvault/vault/record/version1"
	v2 "pvault/vault/record/version2"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type RecordTestSuite struct {
	suite.Suite
	Database version2.Database
	Record   v2.Record
	Password string
}

func TestRecordTestSuite(t *testing.T) {
	suite.Run(t, &RecordTestSuite{})
}

func (s *RecordTestSuite) SetupTest() {
	s.Database = version2.NewDatabase(file.NewPath(s.T(), ""))

	rand := rand.New(0)
	s.Record = v2.NewRecord(rand.ASCII(10))
	s.Record.Username = rand.ASCII(10)
	s.Record.Password = rand.ASCII(30)
	s.Record.Other = map[string]any{"A": rand.ASCII(5), "B": true}

	s.Password = rand.ASCII(30)
}

//=====================================

func (s *RecordTestSuite) TestSaveRecordWithInvalidDatabasePathReturnError() {
	//-- arrange
	s.Database.Path = "invalid/index.bin"

	//-- act
	res := s.Database.SaveRecord(s.Record, s.Password)

	//-- assert
	s.Require().ErrorContains(res, "error writing record file")
}

func (s *RecordTestSuite) TestSaveRecordSavesRecord() {
	//-- act
	res := s.Database.SaveRecord(s.Record, s.Password)

	//-- assert
	s.Require().NoError(res)
	s.Assert().FileExists(s.Database.RecordPath(s.Record.ID))
}

func (s *RecordTestSuite) TestLoadRecordWithInvalidDatabasePathReturnsError() {
	//-- arrange
	err := s.Database.SaveRecord(s.Record, s.Password)
	s.Require().NoError(err)

	s.Database.Path = "invalid/index.bin"

	//-- act
	_, res := s.Database.LoadRecord(s.Record.ID, s.Password)

	//-- assert
	s.Require().ErrorContains(res, "error reading record file")
}

func (s *RecordTestSuite) TestLoadRecordWithIncorrectPasswordReturnsError() {
	//-- arrange
	err := s.Database.SaveRecord(s.Record, s.Password)
	s.Require().NoError(err)

	//-- act
	_, res := s.Database.LoadRecord(s.Record.ID, s.Password+"x")

	//-- assert
	s.Require().ErrorContains(res, "error decrypting ciphertext")
}

func (s *RecordTestSuite) TestLoadRecordWithUnsupportedVersionReturnsError() {
	//-- arrange
	VERSION := version2.RECORD_VERSION + 1

	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, VERSION)

	err := os.WriteFile(s.Database.RecordPath(s.Record.ID), version, 0666)
	s.Require().NoError(err)

	//-- act
	_, res := s.Database.LoadRecord(s.Record.ID, s.Password)

	//-- assert
	s.Require().ErrorContains(res, fmt.Sprintf("unsupported record version \"%d\"", VERSION))
}

func (s *RecordTestSuite) TestLoadRecordVersion1ReturnsRecord() {
	//-- arrange
	s.Record.Other = map[string]any{}
	v1 := v1.Record{
		ID:       s.Record.ID,
		Name:     s.Record.Name,
		Username: s.Record.Username,
		Password: s.Record.Password,
	}

	err := s.Database.SaveRecord(v1, s.Password)
	s.Require().NoError(err)

	//-- act
	r, err := s.Database.LoadRecord(s.Record.ID, s.Password)

	//-- assert
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, r.Upgrade())
}

func (s *RecordTestSuite) TestLoadRecordVersion2ReturnsRecord() {
	//-- arrange
	err := s.Database.SaveRecord(s.Record, s.Password)
	s.Require().NoError(err)

	//-- act
	r, err := s.Database.LoadRecord(s.Record.ID, s.Password)

	//-- assert
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, r)
}

func (s *RecordTestSuite) TestDeleteRecordWithInvalidDatabasePathReturnsError() {
	//-- arrange
	err := s.Database.SaveRecord(s.Record, s.Password)
	s.Require().NoError(err)

	s.Database.Path = "invalid/index.bin"

	//-- act
	res := s.Database.DeleteRecord(s.Record.ID)

	//-- assert
	s.Require().ErrorContains(res, "error removing record file")
}

func (s *RecordTestSuite) TestDeleteRecordReturnsNoError() {
	//-- arrange
	err := s.Database.SaveRecord(s.Record, s.Password)
	s.Require().NoError(err)

	//-- act
	res := s.Database.DeleteRecord(s.Record.ID)

	//-- assert
	s.Require().NoError(res)
	s.Assert().NoFileExists(s.Database.RecordPath(s.Record.ID))
}
