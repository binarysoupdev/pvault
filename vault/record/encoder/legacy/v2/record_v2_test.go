package v2_test

import (
	"bytes"

	v2 "pvault/vault/record/encoder/legacy/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

//TODO: test EncodeDecode
