package v3_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"pvault/util"
	"pvault/vault/index"
	v3 "pvault/vault/index/encoder/v3"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeIndexReturnsErrorWhenErrorWritingHeader(t *testing.T) {
	//-- arrange
	e := v3.Encoder{}
	mock := &util.MockWriter{
		WriteErrors: []error{errors.New("")},
	}

	//-- act
	res := e.EncodeIndex(mock, index.IndexMap{})

	//-- arrange
	assert.ErrorContains(t, res, "error encoding header")
}

func TestEncodeIndexReturnsErrorWhenErrorWritingEntry(t *testing.T) {
	//-- arrange
	e := v3.Encoder{}
	mock := &util.MockWriter{
		WriteErrors: []error{nil, errors.New("")},
	}

	//-- act
	res := e.EncodeIndex(mock, index.IndexMap{"": uuid.Nil})

	//-- arrange
	assert.ErrorContains(t, res, "error encoding entry [0]")
}

func TestDecodeIndexReturnsErrorWhenHeaderTooShort(t *testing.T) {
	//-- arrange
	e := v3.Encoder{}
	buffer := &bytes.Buffer{}

	//-- act
	_, res := e.DecodeIndex(buffer)

	//-- arrange
	assert.ErrorContains(t, res, "error decoding header")
}

func TestDecodeIndexReturnsErrorWhenIndexVersionUnsupported(t *testing.T) {
	//-- arrange
	e := v3.Encoder{}
	buffer := &bytes.Buffer{}

	const VERSION = index.VERSION + 1

	HEADER := make([]byte, 4)
	binary.BigEndian.PutUint16(HEADER, uint16(VERSION))
	buffer.Write(HEADER)

	//-- act
	_, res := e.DecodeIndex(buffer)

	//-- arrange
	assert.ErrorContains(t, res, fmt.Sprintf("unsupported index version \"%d\"", VERSION))
}

func TestDecodeIndexReturnsErrorWhenEntryHeaderTooShort(t *testing.T) {
	//-- arrange
	e := v3.Encoder{}
	buffer := &bytes.Buffer{}

	HEADER := make([]byte, 4)
	binary.BigEndian.PutUint16(HEADER, uint16(index.VERSION))
	binary.BigEndian.PutUint16(HEADER[2:], 1)
	buffer.Write(HEADER)

	//-- act
	_, res := e.DecodeIndex(buffer)

	//-- arrange
	assert.ErrorContains(t, res, "error decoding entry [0]")
	assert.ErrorContains(t, res, "error decoding header")
}

func TestDecodeIndexReturnsErrorWhenEntryLengthTooShort(t *testing.T) {
	//-- arrange
	e := v3.Encoder{}
	buffer := &bytes.Buffer{}

	HEADER := make([]byte, 4)
	binary.BigEndian.PutUint16(HEADER, uint16(index.VERSION))
	binary.BigEndian.PutUint16(HEADER[2:], 1)
	buffer.Write(HEADER)

	const ENTRY_LENGTH = 1

	ENTRY := make([]byte, 2)
	binary.BigEndian.PutUint16(ENTRY, ENTRY_LENGTH)
	buffer.Write(ENTRY)

	//-- act
	_, res := e.DecodeIndex(buffer)

	//-- arrange
	assert.ErrorContains(t, res, "error decoding entry [0]")
	assert.ErrorContains(t, res, fmt.Sprintf("length too short: %d", ENTRY_LENGTH))
}

func TestDecodeIndexReturnsErrorWhenEntryBodyTooShort(t *testing.T) {
	//-- arrange
	e := v3.Encoder{}
	buffer := &bytes.Buffer{}

	HEADER := make([]byte, 4)
	binary.BigEndian.PutUint16(HEADER, uint16(index.VERSION))
	binary.BigEndian.PutUint16(HEADER[2:], 1)
	buffer.Write(HEADER)

	const ENTRY_LENGTH = 16 + 1

	ENTRY := make([]byte, 2)
	binary.BigEndian.PutUint16(ENTRY, uint16(ENTRY_LENGTH))
	buffer.Write(ENTRY)

	//-- act
	_, res := e.DecodeIndex(buffer)

	//-- arrange
	assert.ErrorContains(t, res, "error decoding entry [0]")
	assert.ErrorContains(t, res, "error decoding body")
}

func TestEncodeDecodeIndexReturnsIndexAndNoError(t *testing.T) {
	//-- arrange
	e := v3.Encoder{}
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
