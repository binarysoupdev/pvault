package index_test

import (
	"fmt"
	"pvault/vault/index"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/require"
)

func TestLoadIndexInvalidPath(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")

	//-- act
	_, res := index.LoadIndex(PATH + "/invalid")

	//-- assert
	require.ErrorContains(t, res, "error reading index file")
}

func TestLoadIndexUnsupportedVersion(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	VERSION := 2

	file, PATH := file.Create(t, rand.ASCII(10))
	file.Write([]byte{0, byte(VERSION), 0, 0})
	file.Close()

	//-- act
	_, res := index.LoadIndex(PATH)

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("unsupported version \"%d\"", VERSION))
}
