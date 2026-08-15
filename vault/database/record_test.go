package database_test

import (
	"errors"
	"path/filepath"

	"pvault/vault/database"
	"pvault/vault/record"
	"testing"

	"github.com/binarysoupdev/go-extensions/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveRecordReturnsErrorWhenPathInvalid(t *testing.T) {
	//-- arrange
	PATH := filepath.Join(t.TempDir(), "invalid")
	const PASSWORD = "Password123!"

	RECORD := record.Mock{
		ID: uuid.New(),
	}

	//-- act
	res := database.SaveRecord(&database.EncoderMock{}, PATH, RECORD, PASSWORD)

	//-- assert
	assert.ErrorContains(t, res, "error creating record file")
}

func TestSaveRecordReturnsErrorWhenEncodeRecordReturnsError(t *testing.T) {
	//-- arrange
	const PASSWORD = "Password123!"

	RECORD := record.Mock{
		ID: uuid.New(),
	}

	mock := &database.EncoderMock{
		EncodeRecordError: errors.New(""),
	}

	//-- act
	res := database.SaveRecord(mock, t.TempDir(), RECORD, PASSWORD)

	//-- assert
	assert.ErrorContains(t, res, "error encoding record")
}

func TestSaveRecordReturnsNoErrorAndSavesRecord(t *testing.T) {
	//-- arrange
	PATH := t.TempDir()
	const PASSWORD = "Password123!"

	RECORD := record.Mock{
		ID: uuid.New(),
	}
	mock := &database.EncoderMock{}

	//-- act
	res := database.SaveRecord(mock, PATH, RECORD, PASSWORD)

	//-- assert
	require.NoError(t, res)
	assert.Equal(t, mock.Record, RECORD)
	assert.FileExists(t, mock.RecordPath(PATH, RECORD.ID))
}

func TestLoadRecordReturnsErrorWhenPathInvalid(t *testing.T) {
	//-- act
	_, res := database.LoadRecord(&database.EncoderMock{}, "invalid", uuid.Nil, "")

	//-- assert
	assert.ErrorContains(t, res, "error opening record file")
}

func TestLoadRecordReturnsErrorWhenDecodeRecordReturnsError(t *testing.T) {
	//-- arrange
	mock := &database.EncoderMock{
		DecodeRecordError: errors.New(""),
	}
	ID := uuid.New()

	PATH := t.TempDir()
	require.NoError(t, file.CreateEmpty(mock.RecordPath(PATH, ID)))

	//-- act
	_, res := database.LoadRecord(mock, PATH, ID, "")

	//-- assert
	assert.ErrorContains(t, res, "error decoding record")
}

func TestLoadRecordReturnsRecordAndNoErrorAndLoadsRecord(t *testing.T) {
	//-- arrange
	const PASSWORD = "Password123!"

	mock := &database.EncoderMock{
		Record: record.Mock{
			ID: uuid.New(),
		},
	}

	PATH := t.TempDir()
	require.NoError(t, file.CreateEmpty(mock.RecordPath(PATH, mock.Record.GetID())))

	//-- act
	res, err := database.LoadRecord(mock, PATH, mock.Record.GetID(), PASSWORD)

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, res, mock.Record)
}

func TestDeleteRecordReturnsErrorWhenPathInvalid(t *testing.T) {
	//-- arrange
	mock := &database.EncoderMock{}

	//-- act
	res := database.DeleteRecord(mock, "invalid", uuid.Nil)

	//-- assert
	assert.ErrorContains(t, res, "error removing record file")
}

func TestDeleteRecordReturnsNoErrorAndDeletesRecord(t *testing.T) {
	//-- arrange
	mock := &database.EncoderMock{}
	ID := uuid.New()

	PATH := t.TempDir()
	require.NoError(t, file.CreateEmpty(mock.RecordPath(PATH, ID)))

	//-- act
	res := database.DeleteRecord(mock, PATH, ID)

	//-- assert
	require.NoError(t, res)
	assert.NoFileExists(t, PATH)
}
