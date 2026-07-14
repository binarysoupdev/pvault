package vault_test

import (
	"fmt"
	"path/filepath"
	"pvault/vault"
	"pvault/vault/index/version1"
	"pvault/vault/index/version2"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	VERSION := vault.CURRENT_VERSION + 1

	file, PATH := file.Create(t, version2.INDEX_FILE)
	file.Write([]byte{0, byte(VERSION), 0, 0})
	file.Close()

	//-- act
	_, res := vault.Open(filepath.Dir(PATH))

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("unsupported version \"%d\"", VERSION))
}

func TestOpenLegacyFileLoadsCorrectDecoder(t *testing.T) {
	//-- arrange
	PATH := file.CreateEmpty(t, version1.INDEX_FILE)
	VAULT_PATH := filepath.Dir(PATH)

	//-- act
	v, err := vault.Open(VAULT_PATH)

	//-- assert
	require.NoError(t, err)

	assert.Equal(t, VAULT_PATH, v.Path)
	assert.NotNil(t, v.Index)
	assert.IsType(t, version1.Database{}, v.Database)
}

func TestOpenModernFileLoadsCorrectDecoder(t *testing.T) {
	//-- arrange
	DATABASE := version2.Database{}

	file, PATH := file.Create(t, version2.INDEX_FILE)
	file.Write([]byte{0, byte(DATABASE.GetVersion()), 0, 0})
	file.Close()

	VAULT_PATH := filepath.Dir(PATH)

	//-- act
	v, err := vault.Open(VAULT_PATH)

	//-- assert
	require.NoError(t, err)

	assert.Equal(t, VAULT_PATH, v.Path)
	assert.NotNil(t, v.Index)
	assert.IsType(t, DATABASE, v.Database)
}
