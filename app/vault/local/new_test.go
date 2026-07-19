package local_test

import (
	"path/filepath"
	index_v2 "pvault/app/vault/index/version2"
	vault "pvault/app/vault/local"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreteNewVaultReturnsErrorWhenNewDirectoryAlreadyExists(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")

	//-- act
	_, res := vault.CreateNewVault(PATH)

	//-- assert
	require.ErrorContains(t, res, "vault path already exists")
}

func TestCreteNewVaultReturnsVaultAndNoErrorAndCreatesDirectoryAndSavesMap(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "vault")

	//-- act
	v, res := vault.CreateNewVault(PATH)

	//-- assert
	require.NoError(t, res)

	assert.Equal(t, PATH, v.Path)
	assert.DirExists(t, PATH)
	assert.FileExists(t, filepath.Join(PATH, index_v2.FILENAME))
}

func TestOpenVaultReturnsErrorWhenVaultPathDoesNotContainAndIndexFile(t *testing.T) {
	//-- arrange
	PATH := "invalid"

	//-- act
	_, res := vault.OpenVault(PATH)

	//-- assert
	require.ErrorContains(t, res, "error loading index")
}

func TestOpenVaultReturnsVaultAndNoError(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "vault")
	VAULT, err := vault.CreateNewVault(PATH)
	require.NoError(t, err)

	//-- act
	res, err := vault.OpenVault(PATH)

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, VAULT, res)
}
