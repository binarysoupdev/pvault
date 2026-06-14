package vault_test

import (
	"fmt"
	"pvault/vault"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadIndexInvalidPath(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")

	//-- act
	_, res := vault.LoadIndex(PATH + "/invalid")

	//-- assert
	require.ErrorContains(t, res, "error reading index file")
}

func TestLoadIndexUnsupportedVersion(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	VERSION := 2

	file, PATH := file.Create(t, rand.ASCII(10))
	file.Write([]byte{0, byte(VERSION), 0, 0})
	file.Close()

	//-- act
	_, res := vault.LoadIndex(PATH)

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("unsupported version \"%d\"", VERSION))
}

func TestSaveLoadIndex(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	PATH := file.NewPath(t, rand.ASCII(10))

	INDEX := vault.IndexMap{
		rand.ASCII(10): uuid.New(),
		rand.ASCII(15): uuid.New(),
		rand.ASCII(20): uuid.New(),
	}

	//-- act
	err := INDEX.Save(PATH)
	require.NoError(t, err)

	res, err := vault.LoadIndex(PATH)
	require.NoError(t, err)

	//-- assert
	require.Len(t, res, len(INDEX))

	for key, val := range INDEX {
		assert.Contains(t, res, key)
		assert.Equal(t, val, res[key])
	}
}
