package v2_test

import (
	"bytes"
	"pvault/util"
	"pvault/vault/index"
	v2 "pvault/vault/index/encoder/legacy/v2"
	"testing"

	"github.com/google/uuid"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeIndexReturnsErrorWhenErrorWritingNullHeader(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	mock := &util.MockWriter{
		WriteErrors: []error{errors.New("")},
	}

	//-- act
	res := e.EncodeIndex(mock, index.IndexMap{})

	//-- arrange
	assert.ErrorContains(t, res, "error encoding null header")
}

func TestDecodeIndexReturnsErrorWhenHeaderTooShort(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	buffer := &bytes.Buffer{}

	//-- act
	_, res := e.DecodeIndex(buffer)

	//-- arrange
	assert.ErrorContains(t, res, "error decoding null header")
}

func TestEncodeDecodeIndexReturnsIndexAndNoError(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	buffer := &bytes.Buffer{}

	INDEX := index.IndexMap{
		"":     uuid.Nil,
		"name": uuid.New(),
	}
	require.NoError(t, e.EncodeIndex(buffer, INDEX))

	//-- act
	res, err := e.DecodeIndex(buffer)

	//-- arrange
	require.NoError(t, err)
	assert.Equal(t, INDEX, res)
}
