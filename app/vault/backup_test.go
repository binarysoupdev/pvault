package vault_test

import (
	"fmt"
	"os"
	"path/filepath"
	"pvault/app/vault"
	"pvault/app/vault/database"
	"pvault/app/vault/index"
	"pvault/app/vault/meta"
	"pvault/util"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVaultBackupReturnsErrorWhenPathNotFound(t *testing.T) {
	//-- act
	res := vault.Vault{}.Backup("invalid")

	//-- assert
	assert.ErrorContains(t, res, "error reading backup directory")
}

func TestVaultBackupReturnsErrorWhenPathIsNotADir(t *testing.T) {
	//-- arrange
	PATH := filepath.Join(t.TempDir(), "backups.txt")
	require.NoError(t, util.CreateEmptyFile(PATH))

	//-- act
	res := vault.Vault{}.Backup(PATH)

	//-- assert
	assert.ErrorContains(t, res, fmt.Sprintf("\"%s\" is not a directory", PATH))
}

func TestVaultBackupReturnsErrorWhenErrorBackingMetadataFile(t *testing.T) {
	//-- arrange
	v := vault.Vault{
		Path:            "invalid",
		MetaEncoder:     &meta.EncoderMock{},
		DatabaseEncoder: &database.EncoderMock{},
	}

	//-- act
	res := v.Backup(t.TempDir())

	//-- assert
	assert.ErrorContains(t, res, "error backing metadata file")
}

func TestVaultBackupReturnsErrorWhenErrorBackingIndexFile(t *testing.T) {
	//-- arrange
	v := vault.Vault{
		Path:            t.TempDir(),
		MetaEncoder:     &meta.EncoderMock{},
		DatabaseEncoder: &database.EncoderMock{},
	}

	require.NoError(t, os.WriteFile(v.MetadataPath(), []byte{}, 0666))

	//-- act
	res := v.Backup(t.TempDir())

	//-- assert
	assert.ErrorContains(t, res, "error backing index file")
}

func TestVaultBackupReturnsErrorWhenRecordFilesNotFound(t *testing.T) {
	//-- arrange
	v := vault.Vault{
		Path:            t.TempDir(),
		MetaEncoder:     &meta.EncoderMock{},
		DatabaseEncoder: &database.EncoderMock{},
	}

	R1 := uuid.New()
	R2 := uuid.New()
	v.Map = index.IndexMap{
		"name1": R1,
		"name2": R2,
	}

	require.NoError(t, os.WriteFile(v.MetadataPath(), []byte{}, 0666))
	require.NoError(t, os.WriteFile(v.DatabaseEncoder.IndexPath(v.Path), []byte{}, 0666))

	//-- act
	res := v.Backup(t.TempDir())

	//-- assert
	assert.ErrorContains(t, res, "error backing record "+R1.String())
	assert.ErrorContains(t, res, "error backing record "+R2.String())
}

func TestVaultBackupReturnsNoErrorAndBacksUpIndexAndRecords(t *testing.T) {
	v := vault.Vault{
		Path:            t.TempDir(),
		MetaEncoder:     &meta.EncoderMock{},
		DatabaseEncoder: &database.EncoderMock{},
	}

	R1 := uuid.New()
	R2 := uuid.New()
	v.Map = index.IndexMap{
		"name1": R1,
		"name2": R2,
	}

	METADATA_BYTES := []byte{0, 0, 0}
	require.NoError(t, os.WriteFile(v.MetadataPath(), METADATA_BYTES, 0666))

	INDEX_BYTES := []byte{1, 1, 1}
	require.NoError(t, os.WriteFile(v.IndexPath(), INDEX_BYTES, 0666))

	R1_BYTES := []byte{2, 2, 2}
	require.NoError(t, os.WriteFile(v.RecordPath(R1), R1_BYTES, 0666))

	R2_BYTES := []byte{3, 3, 3}
	require.NoError(t, os.WriteFile(v.RecordPath(R2), R2_BYTES, 0666))

	//-- act
	res := v.Backup(t.TempDir())

	//-- assert
	require.NoError(t, res)

	meta, err := os.ReadFile(v.MetadataPath())
	require.NoError(t, err)
	assert.Equal(t, METADATA_BYTES, meta)

	index, err := os.ReadFile(v.IndexPath())
	require.NoError(t, err)
	assert.Equal(t, INDEX_BYTES, index)

	r1, err := os.ReadFile(v.RecordPath(R1))
	require.NoError(t, err)
	assert.Equal(t, R1_BYTES, r1)

	r2, err := os.ReadFile(v.RecordPath(R2))
	require.NoError(t, err)
	assert.Equal(t, R2_BYTES, r2)
}
