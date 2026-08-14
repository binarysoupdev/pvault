package v1_test

import (
	"bytes"

	v1 "pvault/vault/record/encoder/legacy/v1"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDecodeV1ReturnsErrorWhenErrorDecodingHash(t *testing.T) {
	//-- arrange
	e := v1.Encoder{}
	buffer := &bytes.Buffer{}

	//-- act
	_, res := e.DecodeV1(buffer, "", uuid.Nil, "")

	//-- assert
	assert.ErrorContains(t, res, "error decoding record v1")
	assert.ErrorContains(t, res, "error decoding hash prefix")
}

func TestDecodeV1ReturnsErrorWhenErrorDecryptingRecord(t *testing.T) {
	//-- arrange
	e := v1.Encoder{}
	buffer := &bytes.Buffer{}

	HASH := make([]byte, v1.HASH_SIZE)
	buffer.Write(HASH)

	//-- act
	_, res := e.DecodeV1(buffer, "", uuid.Nil, "")

	//-- assert
	assert.ErrorContains(t, res, "error decrypting record v1")
}

//TODO: test EncodeDecode
