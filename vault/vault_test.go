package vault_test

import (
	"fmt"
	"path/filepath"
	"pvault/vault"
	"pvault/vault/index"
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

func TestOpenIndexFileMissingReturnsError(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")

	//-- act
	_, res := vault.Open(PATH)

	//-- assert
	require.ErrorContains(t, res, "index file not found")
}

func TestOpenUnsupportedVersionReturnsError(t *testing.T) {
	//-- arrange
	VERSION := index.CURRENT_VERSION + 1

	file, PATH := file.Create(t, vault.INDEX_FILE)
	file.Write([]byte{0, byte(VERSION), 0, 0})
	file.Close()

	//-- act
	_, res := vault.Open(filepath.Dir(PATH))

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("unsupported version \"%d\"", VERSION))
}

func TestOpenLegacyFileVersionOutOfDateReturnsError(t *testing.T) {
	//-- arrange
	PATH := file.CreateEmpty(t, vault.LEGACY_INDEX_FILE)

	//-- act
	_, res := vault.Open(filepath.Dir(PATH))

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("version \"%d\" out-of-date", 0))
}

func TestOpenSetsPathAndLoadsIndex(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "vault")

	_, err := vault.InitializeNew(PATH)
	require.NoError(t, err)

	//-- act
	res, err := vault.Open(PATH)

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, PATH, res.Path)
	assert.NotNil(t, res.Index)
}
