package data_test

import (
	"encoding/binary"
	"fmt"
	"os"
	"pvault/vault/data"
	"pvault/vault/record"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/suite"
)

type RecordTestSuite struct {
	suite.Suite
	Database data.DatabaseV2
	Record   record.Record
	Password string
}

func TestRecordTestSuite(t *testing.T) {
	suite.Run(t, &RecordTestSuite{})
}

func (s *RecordTestSuite) SetupTest() {
	s.Database = data.NewDatabaseV2(file.NewPath(s.T(), "index.bin"))

	rand := rand.New(0)
	s.Record = record.NewFromName(rand.ASCII(10))
	s.Record.Username = rand.ASCII(10)
	s.Record.Password = rand.ASCII(30)
	s.Record.Other = []interface{}{rand.ASCII(5), rand.ASCII(5), rand.ASCII(5)}

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
	err := s.Database.SaveRecord(s.Record, s.Password)
	s.Require().NoError(err)

	raw, err := os.ReadFile(s.Database.RecordPath(s.Record.ID))
	s.Require().NoError(err)

	VERSION := data.CURRENT_RECORD_VERSION + 1
	binary.BigEndian.PutUint16(raw, VERSION)

	err = os.WriteFile(s.Database.RecordPath(s.Record.ID), raw, 0666)
	s.Require().NoError(err)

	//-- act
	_, res := s.Database.LoadRecord(s.Record.ID, s.Password)

	//-- assert
	s.Require().ErrorContains(res, fmt.Sprintf("unsupported record version \"%d\"", VERSION))
}

func (s *RecordTestSuite) TestLoadRecordReturnsRecord() {
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
