package v1_test

import (
	"bytes"

	"pvault/vault/index"
	v1 "pvault/vault/index/encoder/legacy/v1"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeIndexReturnsErrorWhenIndexPairInvalid(t *testing.T) {
	//-- arrange
	e := v1.Encoder{}
	buffer := &bytes.Buffer{}

	buffer.Write([]byte("invalid"))

	//-- act
	_, res := e.DecodeIndex(buffer)

	//-- arrange
	assert.ErrorContains(t, res, "[line 1] invalid index pair")
}

func TestDecodeIndexReturnsErrorWhenIDInvalid(t *testing.T) {
	//-- arrange
	e := v1.Encoder{}
	buffer := &bytes.Buffer{}

	buffer.Write([]byte("name:invalid"))

	//-- act
	_, res := e.DecodeIndex(buffer)

	//-- arrange
	assert.ErrorContains(t, res, "[line 1] invalid uuid")
}

func TestEncodeDecodeIndexReturnsIndexAndNoError(t *testing.T) {
	//-- arrange
	e := v1.Encoder{}
	buffer := &bytes.Buffer{}

	INDEX := index.IndexMap{
		"":      uuid.Nil,
		"name1": uuid.New(),
	}
	require.NoError(t, e.EncodeIndex(buffer, INDEX))

	//-- act
	res, err := e.DecodeIndex(buffer)

	//-- arrange
	require.NoError(t, err)
	assert.Equal(t, INDEX, res)
}
