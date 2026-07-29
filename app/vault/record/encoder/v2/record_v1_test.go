package v2_test

import (
	"bytes"
	"encoding/binary"
	v2 "pvault/app/vault/record/encoder/v2"
	record_v1 "pvault/app/vault/record/record/v1"
	"pvault/util"
	"testing"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeV1ReturnsErrorWhenErrorWritingData(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
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
	require.ErrorContains(t, res, "error decrypting record v1")
}
