package database_test

import (
	"errors"
	"path/filepath"
	"pvault/app/vault/database"
	"pvault/app/vault/record"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveRecordReturnsErrorWhenPathInvalid(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "invalid")
	const PASSWORD = "Password123!"

	RECORD := record.Mock{
		ID: uuid.New(),
	}

	//-- act
	res := database.SaveRecord(&database.Mock{}, PATH, RECORD, PASSWORD)

	//-- assert
	assert.ErrorContains(t, res, "error creating record file")
}

func TestSaveRecordReturnsErrorWhenEncodeRecordReturnsError(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")
	const PASSWORD = "Password123!"

	RECORD := record.Mock{
		ID: uuid.New(),
	}

	mock := &database.Mock{
		EncodeRecordError: errors.New(""),
	}

	//-- act
	res := database.SaveRecord(mock, PATH, RECORD, PASSWORD)

	//-- assert
	assert.ErrorContains(t, res, "error encoding record")
}

func TestSaveRecordReturnsNoErrorAndSavesRecord(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")
	const PASSWORD = "Password123!"

	RECORD := record.Mock{
		ID: uuid.New(),
	}
	mock := &database.Mock{}

	//-- act
	res := database.SaveRecord(mock, PATH, RECORD, PASSWORD)

	//-- assert
	require.NoError(t, res)
	assert.Equal(t, mock.Record, RECORD)
	assert.FileExists(t, mock.RecordPath(PATH, RECORD.ID))
}

func TestLoadRecordReturnsErrorWhenPathInvalid(t *testing.T) {
	//-- act
	_, res := database.LoadRecord(&database.Mock{}, "invalid", uuid.Nil, "")

	//-- assert
	assert.ErrorContains(t, res, "error opening record file")
}

func TestLoadRecordReturnsErrorWhenDecodeRecordReturnsError(t *testing.T) {
	//-- arrange
	mock := &database.Mock{
		DecodeRecordError: errors.New(""),
	}

	ID := uuid.New()
	FILEPATH := file.CreateEmpty(t, mock.RecordPath("", ID))

	//-- act
	_, res := database.LoadRecord(mock, filepath.Dir(FILEPATH), ID, "")

	//-- assert
	assert.ErrorContains(t, res, "error decoding record")
}

func TestLoadRecordReturnsRecordAndNoErrorAndLoadsRecord(t *testing.T) {
	//-- arrange
	const PASSWORD = "Password123!"

	mock := &database.Mock{
		Record: record.Mock{
			ID: uuid.New(),
		},
	}
	FILEPATH := file.CreateEmpty(t, mock.RecordPath("", mock.Record.GetID()))

	//-- act
	res, err := database.LoadRecord(mock, filepath.Dir(FILEPATH), mock.Record.GetID(), PASSWORD)

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, res, mock.Record)
}

func TestDeleteRecordReturnsErrorWhenPathInvalid(t *testing.T) {
	//-- arrange
	mock := &database.Mock{}

	//-- act
	res := database.DeleteRecord(mock, "invalid", uuid.Nil)

	//-- assert
	assert.ErrorContains(t, res, "error removing record file")
}

func TestDeleteRecordReturnsNoErrorAndDeletesRecord(t *testing.T) {
	//-- arrange
	mock := &database.Mock{}
	ID := uuid.New()
	FILEPATH := file.CreateEmpty(t, mock.RecordPath("", ID))

	//-- act
	res := database.DeleteRecord(mock, filepath.Dir(FILEPATH), ID)

	//-- assert
	require.NoError(t, res)
	assert.NoFileExists(t, FILEPATH)
}
