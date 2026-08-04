package vault_test

import (
	"path/filepath"
	"pvault/app/vault"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializeNewReturnsErrorWhenNewDirectoryAlreadyExists(t *testing.T) {
	//-- act
	_, res := vault.InitializeNew(t.TempDir(), "")

	//-- assert
	require.ErrorContains(t, res, "vault path already exists")
}

func TestInitializeNewReturnsVaultAndNoErrorAndCreatesDirectoryAndVaultFiles(t *testing.T) {
	//-- arrange
	PATH := filepath.Join(t.TempDir(), "vault")
	const NAME = "nickname"

	//-- act
	res, err := vault.InitializeNew(PATH, NAME)

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, PATH, res.Path)
	assert.Equal(t, NAME, res.Meta.Nickname)

	assert.DirExists(t, PATH)
	assert.FileExists(t, res.IndexPath())
	assert.FileExists(t, res.MetadataPath())
}
