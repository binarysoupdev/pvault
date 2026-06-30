package version2_test

import (
	"fmt"
	"pvault/vault"
	"pvault/vault/data/version2"
	"pvault/vault/index"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveIndexWithInvalidPathReturnsError(t *testing.T) {
	//-- arrange
	db := version2.NewDatabase("invalid/index.bin")

	//-- act
	res := db.SaveIndex(index.IndexMap{})

	//-- assert
	require.ErrorContains(t, res, "error creating index file")
}

func TestSaveIndexValidSavesIndex(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	INDEX_FILE := file.NewPath(t, rand.ASCII(10))

	db := version2.NewDatabase(INDEX_FILE)

	//-- act
	res := db.SaveIndex(index.IndexMap{})

	//-- assert
	require.NoError(t, res)
	assert.FileExists(t, INDEX_FILE)
}

func TestLoadIndexWithFileNotFoundReturnsError(t *testing.T) {
	//-- arrange
	db := version2.NewDatabase("invalid")

	//-- act
	_, res := db.LoadIndex()

	//-- assert
	require.ErrorContains(t, res, "error reading index file")
}

func TestLoadIndexWithIncorrectVersionReturnError(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	VERSION := vault.CURRENT_VERSION + 1

	file, PATH := file.Create(t, rand.ASCII(10))
	file.Write([]byte{0, byte(VERSION), 0, 0})
	file.Close()

	db := version2.NewDatabase(PATH)

	//-- act
	_, res := db.LoadIndex()

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("incorrect version \"%d\"", VERSION))
}

func TestLoadIndexValidReturnIndex(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	db := version2.NewDatabase(file.NewPath(t, rand.ASCII(10)))

	INDEX := index.IndexMap{
		rand.ASCII(10): uuid.New(),
		rand.ASCII(15): uuid.New(),
		rand.ASCII(20): uuid.New(),
	}

	err := db.SaveIndex(INDEX)
	require.NoError(t, err)

	//-- act
	res, err := db.LoadIndex()
	require.NoError(t, err)

	//-- assert
	require.Len(t, res, len(INDEX))

	for key, val := range INDEX {
		assert.Contains(t, res, key)
		assert.Equal(t, val, res[key])
	}
}
