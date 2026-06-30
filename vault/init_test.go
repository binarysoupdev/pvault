package vault_test

import (
	"errors"
	"path/filepath"
	"pvault/vault"
	"pvault/vault/data"
	"pvault/vault/data/version2"
	"pvault/vault/index"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/google/uuid"
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
	assert.FileExists(t, filepath.Join(PATH, version2.INDEX_FILE))
}

func TestReloadIndexWhereDatabaseLoadIndexFailsReturnsError(t *testing.T) {
	//-- arrange
	v := vault.Vault{
		Database: &data.DatabaseMock{
			LoadIndexError: errors.New(""),
		},
	}

	//-- act
	res := v.ReloadIndex()

	//-- assert
	require.ErrorContains(t, res, "error loading index from database")
}

func TestReloadIndexValidUpdatesVaultIndex(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	INDEX := index.IndexMap{
		rand.ASCII(10): uuid.New(),
	}

	v := vault.Vault{
		Database: &data.DatabaseMock{
			Index: INDEX,
		},
	}

	//-- act
	res := v.ReloadIndex()

	//-- assert
	require.NoError(t, res)
	assert.Equal(t, INDEX, v.Index)
}
