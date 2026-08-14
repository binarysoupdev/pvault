package v3_test

import (
	"bytes"
	"encoding/binary"

	v3 "pvault/vault/record/encoder/v3"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDecodeV1ReturnsErrorWhenErrorReadingLength(t *testing.T) {
	//-- arrange
	e := v3.Encoder{}
	buffer := &bytes.Buffer{}

	//-- act
	_, res := e.DecodeV1(buffer, "")

	//-- assert
	assert.ErrorContains(t, res, "error decoding record v1")
	assert.ErrorContains(t, res, "error reading length")
}

func TestDecodeV1ReturnsErrorWhenReadingID(t *testing.T) {
	//-- arrange
	e := v3.Encoder{}
	buffer := &bytes.Buffer{}

	LENGTH := make([]byte, 2)
	buffer.Write(LENGTH)

	//-- act
	_, res := e.DecodeV1(buffer, "")

	//-- assert
	assert.ErrorContains(t, res, "error decoding record v1")
	assert.ErrorContains(t, res, "error reading id")
}

func TestDecodeV1ReturnsErrorWhenErrorReadingName(t *testing.T) {
	//-- arrange
	e := v3.Encoder{}
	buffer := &bytes.Buffer{}

	LENGTH := make([]byte, 2)
	binary.BigEndian.PutUint16(LENGTH, 1)
	buffer.Write(LENGTH)

	buffer.Write(uuid.Nil[:])

	//-- act
	_, res := e.DecodeV1(buffer, "")

	//-- assert
	assert.ErrorContains(t, res, "error decoding record v1")
	assert.ErrorContains(t, res, "error reading name")
}

func TestDecodeV1ReturnsErrorWhenErrorDecryptingRecord(t *testing.T) {
	//-- arrange
	e := v3.Encoder{}
	buffer := &bytes.Buffer{}

	LENGTH := make([]byte, 2)
	buffer.Write(LENGTH)
	buffer.Write(uuid.Nil[:])

	//-- act
	_, res := e.DecodeV1(buffer, "")

	//-- assert
	assert.ErrorContains(t, res, "error decrypting record v1")
}

//TODO: test EncodeDecode
