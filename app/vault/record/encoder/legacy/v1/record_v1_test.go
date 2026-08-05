package v1_test

import (
	"bytes"
	v1 "pvault/app/vault/record/encoder/legacy/v1"
	record_v1 "pvault/app/vault/record/record/legacy/v1"
	"pvault/util"
	"testing"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEncodeV1ReturnsErrorWhenErrorWritingData(t *testing.T) {
	//-- arrange
	e := v1.Encoder{}
	mock := &util.MockWriter{
		WriteErrors: []error{errors.New("")},
	}

	//-- act
	res := e.EncodeV1(mock, "", record_v1.Record{})

	//-- assert
	assert.ErrorContains(t, res, "error encoding record v1")
}

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
