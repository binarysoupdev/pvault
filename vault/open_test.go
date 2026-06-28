package vault_test

import (
	"fmt"
	"path/filepath"
	"pvault/vault"
	"pvault/vault/index"
	"pvault/vault/index/legacy"
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

	file, PATH := file.Create(t, vault.INDEX_FILE)
	file.Write([]byte{0, byte(VERSION), 0, 0})
	file.Close()

	//-- act
	_, res := vault.Open(filepath.Dir(PATH))

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("unsupported version \"%d\"", VERSION))
}

func TestOpenLegacyFileLoadsLegacyDecoder(t *testing.T) {
	//-- arrange
	PATH := file.CreateEmpty(t, vault.LEGACY_INDEX_FILE)

	//-- act
	v, err := vault.Open(filepath.Dir(PATH))

	//-- assert
	require.NoError(t, err)
	assert.IsType(t, legacy.DecoderV1{}, v.Decoder)
}

func TestOpenNewVaultLoadsCurrentDecoder(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "vault")

	_, err := vault.InitializeNew(PATH)
	require.NoError(t, err)

	//-- act
	v, err := vault.Open(PATH)

	//-- assert
	require.NoError(t, err)

	assert.Equal(t, PATH, v.Path)
	assert.Equal(t, vault.CURRENT_VERSION, v.Version)
	assert.IsType(t, index.Codec{}, v.Decoder)
	assert.NotNil(t, v.Index)
}
