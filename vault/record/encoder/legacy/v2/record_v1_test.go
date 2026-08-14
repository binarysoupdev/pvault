package v2_test

import (
	"bytes"
	"encoding/binary"

	v2 "pvault/vault/record/encoder/legacy/v2"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDecodeV1ReturnsErrorWhenErrorReadingLength(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	buffer := &bytes.Buffer{}

	//-- act
	_, res := e.DecodeV1(buffer, "", uuid.Nil)

	//-- assert
	assert.ErrorContains(t, res, "error decoding record v1")
	assert.ErrorContains(t, res, "error reading length")
}

func TestDecodeV1ReturnsErrorWhenErrorReadingName(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	buffer := &bytes.Buffer{}

	LENGTH := make([]byte, 2)
	binary.BigEndian.PutUint16(LENGTH, 1)
	buffer.Write(LENGTH)

	//-- act
	_, res := e.DecodeV1(buffer, "", uuid.Nil)

	//-- assert
	assert.ErrorContains(t, res, "error decoding record v1")
	assert.ErrorContains(t, res, "error reading name")
}

func TestDecodeV1ReturnsErrorWhenErrorDecryptingRecord(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	buffer := &bytes.Buffer{}

	LENGTH := make([]byte, 2)
	buffer.Write(LENGTH)

	//-- act
	_, res := e.DecodeV1(buffer, "", uuid.Nil)

	//-- assert
	assert.ErrorContains(t, res, "error decrypting record v1")
}

//TODO: test EncodeDecode
