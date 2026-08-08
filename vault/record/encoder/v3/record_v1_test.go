package v3_test

import (
	"bytes"
	"encoding/binary"
	"pvault/util"
	v3 "pvault/vault/record/encoder/v3"
	record_v1 "pvault/vault/record/record/legacy/v1"
	"testing"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEncodeV1ReturnsErrorWhenErrorWritingData(t *testing.T) {
	//-- arrange
	e := v3.Encoder{}
	mock := &util.MockWriter{
		WriteErrors: []error{errors.New("")},
	}

	RECORD := record_v1.Record{}

	//-- act
	res := e.EncodeV1(mock, "", RECORD)

	//-- assert
	assert.ErrorContains(t, res, "error encoding record v1")
}

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
