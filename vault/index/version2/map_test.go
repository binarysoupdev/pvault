package v2_test

import (
	"fmt"
	"path/filepath"
	"pvault/vault/data"
	v2 "pvault/vault/index/version2"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveMapReturnsErrorWhenPathInvalid(t *testing.T) {
	//-- arrange
	idx := v2.NewIndex("invalid/index.bin")

	//-- act
	res := idx.SaveMap(data.NameMap{})

	//-- assert
	require.ErrorContains(t, res, "error creating index file")
}

func TestSaveMapReturnsNoErrorAndSavesMap(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")

	idx := v2.NewIndex(PATH)

	//-- act
	res := idx.SaveMap(data.NameMap{})

	//-- assert
	require.NoError(t, res)
	assert.FileExists(t, idx.Filepath())
}

func TestLoadMapReturnsErrorWhenFileNotFound(t *testing.T) {
	//-- arrange
	idx := v2.NewIndex("invalid")

	//-- act
	_, res := idx.LoadMap()

	//-- assert
	require.ErrorContains(t, res, "error reading index file")
}

func TestLoadMapReturnsErrorWhenVersionIncorrect(t *testing.T) {
	//-- arrange
	VERSION := v2.VERSION + 1

	file, PATH := file.Create(t, v2.FILENAME)
	file.Write([]byte{0, byte(VERSION), 0, 0})
	file.Close()

	idx := v2.NewIndex(filepath.Dir(PATH))

	//-- act
	_, res := idx.LoadMap()

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("incorrect version \"%d\"", VERSION))
}

func TestLoadMapReturnsMapAndNoError(t *testing.T) {
	//-- arrange
	idx := v2.NewIndex(file.NewPath(t, ""))

	MAP := data.NameMap{
		"name1": uuid.New(),
		"name2": uuid.New(),
		"name3": uuid.New(),
	}

	err := idx.SaveMap(MAP)
	require.NoError(t, err)

	//-- act
	res, err := idx.LoadMap()
	require.NoError(t, err)

	//-- assert
	require.Len(t, res, len(MAP))

	for key, val := range MAP {
		assert.Contains(t, res, key)
		assert.Equal(t, val, res[key])
	}
}
