package vault_test

import (
	"fmt"
	"os"
	"pvault/vault"
	"pvault/vault/database"
	db_v1 "pvault/vault/database/encoder/legacy/v1"
	db_v2 "pvault/vault/database/encoder/legacy/v2"
	db_v3 "pvault/vault/database/encoder/v3"
	"pvault/vault/index"
	"pvault/vault/meta"
	meta_v1 "pvault/vault/meta/encoder/v1"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenReturnsErrorWhenDatabaseNotFound(t *testing.T) {
	//-- arrange
	PATH := t.TempDir()

	//-- act
	_, res := vault.Open(PATH)

	//-- assert
	assert.ErrorContains(t, res, fmt.Sprintf("vault not found at \"%s\"", PATH))
}

func TestOpenReturnsVaultAndNoErrorWhenDatabaseIsV1AndCreatesNewMetadata(t *testing.T) {
	//-- arrange
	PATH := t.TempDir()

	INDEX := index.IndexMap{
		"name": uuid.New(),
	}
	require.NoError(t, database.SaveIndex(db_v1.Encoder{}, PATH, INDEX))

	//-- act
	res, err := vault.Open(PATH)

	//-- assert
	require.NoError(t, err)
	assert.IsType(t, db_v1.Encoder{}, res.DatabaseEncoder)

	assert.Equal(t, PATH, res.Path)
	assert.Equal(t, INDEX, res.Map)

	assert.FileExists(t, res.MetadataPath())
	assert.Equal(t, db_v1.VERSION, res.Meta.DatabaseVersion)
}

func TestOpenReturnsVaultAndNoErrorWhenDatabaseIsV2AndCreatesNewMetadata(t *testing.T) {
	//-- arrange
	PATH := t.TempDir()

	INDEX := index.IndexMap{
		"name": uuid.New(),
	}
	require.NoError(t, database.SaveIndex(db_v2.Encoder{}, PATH, INDEX))

	//-- act
	res, err := vault.Open(PATH)

	//-- assert
	require.NoError(t, err)
	assert.IsType(t, db_v2.Encoder{}, res.DatabaseEncoder)

	assert.Equal(t, PATH, res.Path)
	assert.Equal(t, INDEX, res.Map)

	assert.FileExists(t, res.MetadataPath())
	assert.Equal(t, db_v2.VERSION, res.Meta.DatabaseVersion)
}

func TestOpenReturnsErrorWhenErrorLoadingMetadata(t *testing.T) {
	//-- arrange
	PATH := t.TempDir()
	require.NoError(t, os.WriteFile(meta_v1.Encoder{}.MetadataPath(PATH), []byte{}, 0666))

	//-- act
	_, res := vault.Open(PATH)

	//-- assert
	assert.ErrorContains(t, res, "error loading metadata")
}

func TestOpenReturnsVaultAndNoErrorWhenDatabaseIsV3AndLoadsMetadata(t *testing.T) {
	//-- arrange
	PATH := t.TempDir()
	INDEX := index.IndexMap{
		"name": uuid.New(),
	}
	require.NoError(t, database.SaveIndex(db_v3.Encoder{}, PATH, INDEX))

	VAULT := vault.Vault{
		Path:        PATH,
		MetaEncoder: meta_v1.Encoder{},
		Meta: meta.Metadata{
			DatabaseVersion: db_v3.VERSION,
		},
	}
	require.NoError(t, VAULT.SaveMetadata())

	//-- act
	res, err := vault.Open(PATH)

	//-- assert
	require.NoError(t, err)
	assert.IsType(t, db_v3.Encoder{}, res.DatabaseEncoder)

	assert.Equal(t, PATH, res.Path)
	assert.Equal(t, INDEX, res.Map)
	assert.Equal(t, VAULT.Meta, res.Meta)
}

func TestOpenReturnsErrorWhenDatabaseVersionUnsupported(t *testing.T) {
	//-- arrange
	PATH := t.TempDir()

	VAULT := vault.Vault{
		Path:        PATH,
		MetaEncoder: meta_v1.Encoder{},
		Meta: meta.Metadata{
			DatabaseVersion: db_v3.VERSION + 1,
		},
	}
	require.NoError(t, VAULT.SaveMetadata())

	//-- act
	_, res := vault.Open(PATH)

	//-- assert
	assert.ErrorContains(t, res, fmt.Sprintf("unsupported vault version \"%d\"", VAULT.Meta.DatabaseVersion))
}
