package database_test

import (
	"errors"
	"path/filepath"
	"pvault/app/vault/database"
	"pvault/app/vault/index"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveIndexReturnsErrorWhenPathInvalid(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "invalid")

	//-- act
	res := database.SaveIndex(&database.Mock{}, PATH, index.IndexMap{})

	//-- assert
	assert.ErrorContains(t, res, "error creating index file")
}

func TestSaveIndexReturnsErrorWhenEncodeIndexReturnsError(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")
	mock := &database.Mock{
		EncodeIndexError: errors.New(""),
	}

	//-- act
	res := database.SaveIndex(mock, PATH, index.IndexMap{})

	//-- assert
	assert.ErrorContains(t, res, "error encoding index")
}

func TestSaveIndexReturnsNoErrorAndSavesIndex(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")
	INDEX := index.IndexMap{
		"name": uuid.New(),
	}

	mock := &database.Mock{}

	//-- act
	res := database.SaveIndex(mock, PATH, INDEX)

	//-- assert
	require.NoError(t, res)
	assert.Equal(t, mock.Index, INDEX)
	assert.FileExists(t, mock.IndexPath(PATH))
}

func TestLoadIndexReturnsErrorWhenPathInvalid(t *testing.T) {
	//-- act
	_, res := database.LoadIndex(&database.Mock{}, "invalid")

	//-- assert
	assert.ErrorContains(t, res, "error opening index file")
}

func TestLoadIndexReturnsErrorWhenDecodeIndexReturnsError(t *testing.T) {
	//-- arrange
	mock := &database.Mock{
		DecodeIndexError: errors.New(""),
	}
	FILEPATH := file.CreateEmpty(t, mock.IndexPath(""))

	//-- act
	_, res := database.LoadIndex(mock, filepath.Dir(FILEPATH))

	//-- assert
	assert.ErrorContains(t, res, "error decoding index")
}

func TestLoadIndexReturnsIndexAndNoErrorAndLoadsIndex(t *testing.T) {
	//-- arrange
	mock := &database.Mock{
		Index: index.IndexMap{
			"name": uuid.New(),
		},
	}
	FILEPATH := file.CreateEmpty(t, mock.IndexPath(""))

	//-- act
	res, err := database.LoadIndex(mock, filepath.Dir(FILEPATH))

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, res, mock.Index)
}
