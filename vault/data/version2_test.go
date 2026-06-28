package data_test

import (
	"fmt"
	"pvault/vault"
	"pvault/vault/data"
	"pvault/vault/index"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadIndexFileNotFoundReturnsError(t *testing.T) {
	//-- act
	_, res := data.NewDatabaseV2("invalid").LoadIndex()

	//-- assert
	require.ErrorContains(t, res, "error reading index file")
}

func TestLoadIndexIncorrectVersionReturnError(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	VERSION := vault.CURRENT_VERSION + 1

	file, PATH := file.Create(t, rand.ASCII(10))
	file.Write([]byte{0, byte(VERSION), 0, 0})
	file.Close()

	//-- act
	_, res := data.NewDatabaseV2(PATH).LoadIndex()

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("incorrect version \"%d\"", VERSION))
}

func TestSaveIndexThenLoadIndex(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	PATH := file.NewPath(t, rand.ASCII(10))

	INDEX := index.IndexMap{
		rand.ASCII(10): uuid.New(),
		rand.ASCII(15): uuid.New(),
		rand.ASCII(20): uuid.New(),
	}

	db := data.NewDatabaseV2(PATH)

	//-- act
	err := db.SaveIndex(INDEX)
	require.NoError(t, err)

	res, err := db.LoadIndex()
	require.NoError(t, err)

	//-- assert
	require.Len(t, res, len(INDEX))

	for key, val := range INDEX {
		assert.Contains(t, res, key)
		assert.Equal(t, val, res[key])
	}
}
