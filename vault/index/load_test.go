package index_test

import (
	"fmt"
	"path/filepath"
	"pvault/vault/index"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/stretchr/testify/require"
)

func TestLoadIndexIndexFileMissing(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")

	//-- act
	_, res := index.LoadIndex(PATH + "/invalid")

	//-- assert
	require.ErrorContains(t, res, "index file not found")
}

func TestLoadIndexUnsupportedVersion(t *testing.T) {
	//-- arrange
	VERSION := 2

	file, PATH := file.Create(t, index.INDEX_FILE)
	file.Write([]byte{0, byte(VERSION), 0, 0})
	file.Close()

	//-- act
	_, res := index.LoadIndex(filepath.Dir(PATH))

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("unsupported version \"%d\"", VERSION))
}

func TestLoadIndexLegacyFileVersionOutOfDate(t *testing.T) {
	//-- arrange
	PATH := file.CreateEmpty(t, index.LEGACY_INDEX_FILE)

	//-- act
	_, res := index.LoadIndex(filepath.Dir(PATH))

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("version \"%d\" out-of-date", 0))
}
