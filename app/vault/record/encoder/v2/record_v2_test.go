package v2_test

import (
	"bytes"
	v2 "pvault/app/vault/record/encoder/v2"
	record_v2 "pvault/app/vault/record/record/v2"
	"pvault/util"
	"testing"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/stretchr/testify/assert"
)

func TestEncodeV2ReturnsErrorWhenErrorWritingData(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	mock := &util.MockWriter{
		WriteErrors: []error{errors.New("")},
	}

	RECORD := record_v2.Record{}

	//-- act
	res := e.EncodeV2(mock, "", RECORD)

	//-- assert
	assert.ErrorContains(t, res, "error encoding record v2")
}

func TestDecodeV2ReturnsErrorWhenErrorReadingData(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	mock := &util.MockReader{
		ReadErrors: []error{errors.New("")},
	}

	//-- act
	_, res := e.DecodeV2(mock, "")

	//-- assert
	assert.ErrorContains(t, res, "error decoding record v2")
}

func TestDecodeV2ReturnsErrorWhenErrorDecryptingRecord(t *testing.T) {
	//-- arrange
	e := v2.Encoder{}
	buffer := &bytes.Buffer{}

	LENGTH := make([]byte, 2)
	buffer.Write(LENGTH)

	//-- act
	_, res := e.DecodeV2(buffer, "")

	//-- assert
	assert.ErrorContains(t, res, "error decrypting record v2")
}
