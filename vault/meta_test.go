package vault_test

import (
	"os"
	"pvault/vault"
	"pvault/vault/meta"
	"testing"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVaultSaveMetadataReturnsErrorWhenMetaEncoderEncodeMetadataReturnsError(t *testing.T) {
	//-- arrange
	v := vault.Vault{
		Path: t.TempDir(),
		MetaEncoder: &meta.EncoderMock{
			EncodeMetadataError: errors.New(""),
		},
	}

	//-- act
	res := v.SaveMetadata()

	//-- assert
	assert.ErrorContains(t, res, "error saving metadata")
}

func TestVaultSaveMetadataReturnsNoErrorAndSavesMetadata(t *testing.T) {
	//-- arrange
	mock := &meta.EncoderMock{}

	v := vault.Vault{
		Path:        t.TempDir(),
		MetaEncoder: mock,
		Meta: meta.Metadata{
			DatabaseVersion: 1,
			Nickname:        "nickname",
		},
	}

	//-- act
	res := v.SaveMetadata()

	//-- assert
	require.NoError(t, res)
	assert.FileExists(t, v.MetadataPath())
	assert.Equal(t, v.Meta, mock.Metadata)
}

func TestVaultLoadMetadataReturnsErrorWhenMetaEncoderDecodeMetadataReturnsError(t *testing.T) {
	//-- arrange
	v := vault.Vault{
		Path: t.TempDir(),
		MetaEncoder: &meta.EncoderMock{
			DecodeMetadataError: errors.New(""),
		},
	}

	//-- act
	res := v.LoadMetadata()

	//-- assert
	assert.ErrorContains(t, res, "error loading metadata")
}

func TestVaultLoadMetadataReturnsNoErrorAndLoadsMetadata(t *testing.T) {
	//-- arrange
	mock := &meta.EncoderMock{
		Metadata: meta.Metadata{
			DatabaseVersion: 1,
			Nickname:        "nickname",
		},
	}

	v := vault.Vault{
		Path:        t.TempDir(),
		MetaEncoder: mock,
	}
	require.NoError(t, os.WriteFile(v.MetadataPath(), []byte{}, 0666))

	//-- act
	res := v.LoadMetadata()

	//-- assert
	require.NoError(t, res)
	assert.Equal(t, mock.Metadata, v.Meta)
}
