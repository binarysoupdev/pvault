package vault_test

import (
	"path/filepath"
	"pvault/vault"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializeNewDirectoryExists(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")

	//-- act
	res := vault.InitializeNew(PATH)

	//-- assert
	require.ErrorContains(t, res, "vault path already exists")
}

func TestInitializeNewSuccess(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "vault")

	//-- act
	res := vault.InitializeNew(PATH)

	//-- assert
	require.NoError(t, res)

	assert.DirExists(t, PATH)
	assert.FileExists(t, filepath.Join(PATH, vault.INDEX_FILE))
}

func TestOpenLoadIndexError(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "vault")

	//-- act
	_, res := vault.Open(PATH)

	//-- assert
	require.ErrorContains(t, res, "error loading index file")
}

func TestOpenSuccess(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "vault")

	err := vault.InitializeNew(PATH)
	require.NoError(t, err)

	//-- act
	res, err := vault.Open(PATH)

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, PATH, res.Path)
	assert.NotNil(t, res.Index)
}
