package v1_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"pvault/util"
	"pvault/vault/meta"
	v1 "pvault/vault/meta/encoder/v1"
	"testing"
	"time"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeMetadataReturnsErrorWhenErrorWritingData(t *testing.T) {
	//-- arrange
	const ERROR = "write error"

	e := v1.Encoder{}
	mock := &util.MockWriter{
		WriteErrors: []error{errors.New(ERROR)},
	}

	//-- act
	res := e.EncodeMetadata(mock, meta.Metadata{})

	//-- assert
	assert.ErrorContains(t, res, ERROR)
}

func TestDecodeMetadataReturnsErrorWhenErrorReadingHeader(t *testing.T) {
	//-- arrange
	e := v1.Encoder{}
	buffer := &bytes.Buffer{}

	//-- act
	_, res := e.DecodeMetadata(buffer)

	//-- assert
	assert.ErrorContains(t, res, "error reading header")
}

func TestDecodeMetadataReturnsErrorWhenVersionNotSupported(t *testing.T) {
	//-- arrange
	e := v1.Encoder{}
	buffer := &bytes.Buffer{}

	const VERSION = meta.VERSION + 1

	HEADER := make([]byte, v1.HEADER_SIZE)
	binary.BigEndian.PutUint16(HEADER, VERSION)
	buffer.Write(HEADER)

	//-- act
	_, res := e.DecodeMetadata(buffer)

	//-- assert
	assert.ErrorContains(t, res, fmt.Sprintf("unsupported metadata version \"%d\"", VERSION))
}

func TestDecodeMetadataReturnsErrorWhenErrorReadingNickname(t *testing.T) {
	e := v1.Encoder{}
	buffer := &bytes.Buffer{}

	HEADER := make([]byte, v1.HEADER_SIZE)
	binary.BigEndian.PutUint16(HEADER, meta.VERSION)
	binary.BigEndian.PutUint16(HEADER[2+2:], 1)
	buffer.Write(HEADER)

	//-- act
	_, res := e.DecodeMetadata(buffer)

	//-- assert
	assert.ErrorContains(t, res, "error reading nickname")
}

func TestDecodeMetadataReturnsErrorWhenErrorReadingCreationDate(t *testing.T) {
	e := v1.Encoder{}
	buffer := &bytes.Buffer{}

	HEADER := make([]byte, v1.HEADER_SIZE)
	binary.BigEndian.PutUint16(HEADER, meta.VERSION)
	binary.BigEndian.PutUint16(HEADER[2+2+2:], 1)
	buffer.Write(HEADER)

	//-- act
	_, res := e.DecodeMetadata(buffer)

	//-- assert
	assert.ErrorContains(t, res, "error reading creation date")
}

func TestDecodeMetadataReturnsErrorWhenErrorUnmarshalingCreationDate(t *testing.T) {
	e := v1.Encoder{}
	buffer := &bytes.Buffer{}

	HEADER := make([]byte, v1.HEADER_SIZE)
	binary.BigEndian.PutUint16(HEADER, meta.VERSION)
	binary.BigEndian.PutUint16(HEADER[2+2+2:], 1)
	buffer.Write(HEADER)

	buffer.Write([]byte{0})

	//-- act
	_, res := e.DecodeMetadata(buffer)

	//-- assert
	assert.ErrorContains(t, res, "error unmarshaling creation date")
}

func TestEncodeDecodeReturnsMetadataAndNoError(t *testing.T) {
	e := v1.Encoder{}
	buffer := &bytes.Buffer{}

	DATE := time.Now()
	META := meta.Metadata{
		DatabaseVersion: 1,
		Nickname:        "nickname",
		CreationDate:    DATE,
	}
	require.NoError(t, e.EncodeMetadata(buffer, META))

	//-- act
	res, err := e.DecodeMetadata(buffer)

	//-- assert
	require.NoError(t, err)

	date := res.CreationDate
	res.CreationDate = time.Time{}
	META.CreationDate = time.Time{}

	assert.Equal(t, META, res)
	assert.True(t, date.Equal(DATE))
}
