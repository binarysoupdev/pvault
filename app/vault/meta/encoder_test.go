package meta_test

import (
	"os"
	"path/filepath"
	"pvault/app/vault/meta"
	"testing"
	"time"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveMetadataReturnsErrorWhenPathInvalid(t *testing.T) {
	//-- arrange
	PATH := filepath.Join(t.TempDir(), "path/invalid")

	//-- act
	res := meta.SaveMetadata(&meta.EncoderMock{}, PATH, meta.Metadata{})

	//-- assert
	assert.ErrorContains(t, res, "error creating metadata file")
}

func TestSaveMetadataReturnsErrorWhenEncodeMetadataReturnsError(t *testing.T) {
	//-- arrange
	mock := &meta.EncoderMock{
		EncodeMetadataError: errors.New(""),
	}

	//-- act
	res := meta.SaveMetadata(mock, t.TempDir(), meta.Metadata{})

	//-- assert
	assert.ErrorContains(t, res, "error encoding metadata")
}

func TestSaveMetadataReturnsNoErrorAndSavesMetadata(t *testing.T) {
	//-- arrange
	PATH := t.TempDir()

	META := meta.Metadata{
		DatabaseVersion: 1,
		Nickname:        "nickname",
		CreationDate:    time.Now(),
	}
	mock := &meta.EncoderMock{}

	//-- act
	res := meta.SaveMetadata(mock, PATH, META)

	//-- assert
	require.NoError(t, res)
	assert.Equal(t, mock.Metadata, META)

	assert.FileExists(t, mock.MetadataPath(PATH))
}

func TestLoadMetadataReturnsErrorWhenPathInvalid(t *testing.T) {
	//-- act
	_, res := meta.LoadMetadata(&meta.EncoderMock{}, "invalid")

	//-- assert
	assert.ErrorContains(t, res, "error opening metadata file")
}

func TestLoadMetadataReturnsErrorWhenDecodeMetadataReturnsError(t *testing.T) {
	//-- arrange
	PATH := t.TempDir()

	mock := &meta.EncoderMock{
		DecodeMetadataError: errors.New(""),
	}
	require.NoError(t, os.WriteFile(mock.MetadataPath(PATH), []byte{}, 0666))

	//-- act
	_, res := meta.LoadMetadata(mock, PATH)

	//-- assert
	assert.ErrorContains(t, res, "error decoding metadata")
}

func TestLoadMetadataReturnsMetadataAndNoErrorAndLoadsMetadata(t *testing.T) {
	//-- arrange
	PATH := t.TempDir()

	mock := &meta.EncoderMock{
		Metadata: meta.Metadata{
			DatabaseVersion: 1,
			Nickname:        "nickname",
			CreationDate:    time.Now(),
		},
	}
	require.NoError(t, os.WriteFile(mock.MetadataPath(PATH), []byte{}, 0666))

	//-- act
	res, err := meta.LoadMetadata(mock, PATH)

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, res, mock.Metadata)
}
