package v1_test

import (
	"bytes"
	"pvault/app/vault/index"
	v1 "pvault/app/vault/index/encoder/v1"
	"pvault/util"
	"testing"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeIndexReturnsErrorWhenErrorWritingEntry(t *testing.T) {
	//-- arrange
	e := v1.Encoder{}
	mock := &util.MockWriter{
		WriteErrors: []error{errors.New("")},
	}

	//-- act
	res := e.EncodeIndex(mock, index.IndexMap{"": uuid.Nil})

	//-- arrange
	assert.ErrorContains(t, res, "error encoding entry [0]")
}

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

	err := e.EncodeIndex(buffer, INDEX)
	require.NoError(t, err)

	//-- act
	res, err := e.DecodeIndex(buffer)

	//-- arrange
	require.NoError(t, err)
	assert.Equal(t, INDEX, res)
}
