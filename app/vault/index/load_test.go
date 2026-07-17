package index_test

import (
	"fmt"
	"path/filepath"
	"pvault/app/vault/data"
	"pvault/app/vault/index"
	v1 "pvault/app/vault/index/version1"
	v2 "pvault/app/vault/index/version2"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadReturnsErrorWhenFileNotFound(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")

	//-- act
	_, res := index.Load(PATH)

	//-- assert
	require.ErrorContains(t, res, "index file not found")
}

func TestLoadReturnsErrorWhenVersionUnsupported(t *testing.T) {
	//-- arrange
	VERSION := v2.VERSION + 1

	file, PATH := file.Create(t, v2.FILENAME)
	file.Write([]byte{0, byte(VERSION), 0, 0})
	file.Close()

	//-- act
	_, res := index.Load(filepath.Dir(PATH))

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("unsupported index version \"%d\"", VERSION))
}

func TestLoadReturnsV2IndexAndNoErrorWhenV2Detected(t *testing.T) {
	//-- arrange
	INDEX := v2.NewIndex(file.NewPath(t, ""))
	require.NoError(t, INDEX.SaveMap(data.NameMap{}))

	//-- act
	idx, err := index.Load(INDEX.Path)

	//-- assert
	require.NoError(t, err)
	assert.IsType(t, INDEX, idx)
}

func TestLoadReturnsV1IndexAndNoErrorWhenLegacyFileDetected(t *testing.T) {
	//-- arrange
	INDEX := v1.NewIndex(file.NewPath(t, ""))
	require.NoError(t, INDEX.SaveMap(data.NameMap{}))

	//-- act
	idx, err := index.Load(INDEX.Path)

	//-- assert
	require.NoError(t, err)
	assert.IsType(t, INDEX, idx)
}
