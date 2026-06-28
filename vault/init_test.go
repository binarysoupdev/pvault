package vault_test

import (
	"path/filepath"
	"pvault/vault"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializeNewDirectoryExistsReturnsError(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")

	//-- act
	_, res := vault.InitializeNew(PATH)

	//-- assert
	require.ErrorContains(t, res, "vault path already exists")
}

func TestInitializeNewCreatesDirectoryAndIndexFile(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "vault")

	//-- act
	_, res := vault.InitializeNew(PATH)

	//-- assert
	require.NoError(t, res)

	assert.DirExists(t, PATH)
	assert.FileExists(t, filepath.Join(PATH, vault.INDEX_FILE))
}
