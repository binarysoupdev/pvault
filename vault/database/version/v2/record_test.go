package v2_test

import (
	"encoding/binary"
	"fmt"
	"os"
	v1 "pvault/vault/database/version/v1"
	v2 "pvault/vault/database/version/v2"
	"pvault/vault/record"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type RecordTestSuite struct {
	suite.Suite
	Database v2.Database
	Record   record.RecordV2
	Password string
}

func TestRecordTestSuite(t *testing.T) {
	suite.Run(t, &RecordTestSuite{})
}

func (s *RecordTestSuite) SetupTest() {
	s.Database = v2.New(file.NewPath(s.T(), ""))

	rand := rand.New(0)
	s.Record = record.NewFromName(rand.ASCII(10))
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
	s.Require().ErrorContains(res, "error creating record file")
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
	VERSION := v2.RECORD_VERSION + 1

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

	err := v1.New(s.Database.Path).SaveRecordV1(s.Database.RecordPath(s.Record.ID), s.Record, s.Password)
	s.Require().NoError(err)

	//-- act
	r, err := s.Database.LoadRecord(s.Record.ID, s.Password)

	//-- assert
	s.Require().NoError(err)
	s.Assert().Equal(s.Record, r)
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
