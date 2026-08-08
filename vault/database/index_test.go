package database_test

import (
	"errors"
	"path/filepath"
	"pvault/util"
	"pvault/vault/database"
	"pvault/vault/index"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveIndexReturnsErrorWhenPathInvalid(t *testing.T) {
	//-- arrange
	PATH := filepath.Join(t.TempDir(), "invalid")

	//-- act
	res := database.SaveIndex(&database.EncoderMock{}, PATH, index.IndexMap{})

	//-- assert
	assert.ErrorContains(t, res, "error creating index file")
}

func TestSaveIndexReturnsErrorWhenEncodeIndexReturnsError(t *testing.T) {
	//-- arrange
	mock := &database.EncoderMock{
		EncodeIndexError: errors.New(""),
	}

	//-- act
	res := database.SaveIndex(mock, t.TempDir(), index.IndexMap{})

	//-- assert
	assert.ErrorContains(t, res, "error encoding index")
}

func TestSaveIndexReturnsNoErrorAndSavesIndex(t *testing.T) {
	//-- arrange
	PATH := t.TempDir()
	INDEX := index.IndexMap{
		"name": uuid.New(),
	}

	mock := &database.EncoderMock{}

	//-- act
	res := database.SaveIndex(mock, PATH, INDEX)

	//-- assert
	require.NoError(t, res)
	assert.Equal(t, mock.Index, INDEX)
	assert.FileExists(t, mock.IndexPath(PATH))
}

func TestLoadIndexReturnsErrorWhenPathInvalid(t *testing.T) {
	//-- act
	_, res := database.LoadIndex(&database.EncoderMock{}, "invalid")

	//-- assert
	assert.ErrorContains(t, res, "error opening index file")
}

func TestLoadIndexReturnsErrorWhenDecodeIndexReturnsError(t *testing.T) {
	//-- arrange
	mock := &database.EncoderMock{
		DecodeIndexError: errors.New(""),
	}

	PATH := t.TempDir()
	require.NoError(t, util.CreateEmptyFile(mock.IndexPath(PATH)))

	//-- act
	_, res := database.LoadIndex(mock, PATH)

	//-- assert
	assert.ErrorContains(t, res, "error decoding index")
}

func TestLoadIndexReturnsIndexAndNoErrorAndLoadsIndex(t *testing.T) {
	//-- arrange
	mock := &database.EncoderMock{
		Index: index.IndexMap{
			"name": uuid.New(),
		},
	}

	PATH := t.TempDir()
	require.NoError(t, util.CreateEmptyFile(mock.IndexPath(PATH)))

	//-- act
	res, err := database.LoadIndex(mock, PATH)

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, res, mock.Index)
}
