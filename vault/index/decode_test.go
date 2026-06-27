package index_test

import (
	"fmt"
	"pvault/vault/index"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/require"
)

func TestDecodeIndexFileNotFoundReturnsError(t *testing.T) {
	//-- act
	_, res := index.Codec{}.Decode("invalid")

	//-- assert
	require.ErrorContains(t, res, "error reading index file")
}

func TestDecodeIncorrectVersionReturnError(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	VERSION := index.CURRENT_VERSION + 1

	file, PATH := file.Create(t, rand.ASCII(10))
	file.Write([]byte{0, byte(VERSION), 0, 0})
	file.Close()

	//-- act
	_, res := index.Codec{}.Decode(PATH)

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("incorrect version \"%d\"", VERSION))
}
