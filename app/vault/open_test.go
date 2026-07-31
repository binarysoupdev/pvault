package vault_test

import (
	"fmt"
	"os"
	"pvault/app/vault"
	"pvault/app/vault/database"
	db_v1 "pvault/app/vault/database/database/v1"
	db_v2 "pvault/app/vault/database/database/v2"
	db_v3 "pvault/app/vault/database/database/v3"
	"pvault/app/vault/index"
	"pvault/app/vault/meta"
	meta_encoder "pvault/app/vault/meta/encoder"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenReturnsErrorWhenDatabaseNotFound(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")

	//-- act
	_, res := vault.Open(PATH)

	//-- assert
	assert.ErrorContains(t, res, fmt.Sprintf("vault not found at \"%s\"", PATH))
}

func TestOpenReturnsVaultAndNoErrorWhenDatabaseIsV1AndCreatesNewMetadata(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")

	INDEX := index.IndexMap{
		"name": uuid.New(),
	}
	require.NoError(t, database.SaveIndex(db_v1.Database{}, PATH, INDEX))

	//-- act
	res, err := vault.Open(PATH)

	//-- assert
	require.NoError(t, err)
	assert.IsType(t, db_v1.Database{}, res.Database)

	assert.Equal(t, PATH, res.Path)
	assert.Equal(t, INDEX, res.Map)

	assert.FileExists(t, res.MetadataPath())
	assert.Equal(t, db_v1.VERSION, res.Meta.DatabaseVersion)
}

func TestOpenReturnsVaultAndNoErrorWhenDatabaseIsV2AndCreatesNewMetadata(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")

	INDEX := index.IndexMap{
		"name": uuid.New(),
	}
	require.NoError(t, database.SaveIndex(db_v2.Database{}, PATH, INDEX))

	//-- act
	res, err := vault.Open(PATH)

	//-- assert
	require.NoError(t, err)
	assert.IsType(t, db_v2.Database{}, res.Database)

	assert.Equal(t, PATH, res.Path)
	assert.Equal(t, INDEX, res.Map)

	assert.FileExists(t, res.MetadataPath())
	assert.Equal(t, db_v2.VERSION, res.Meta.DatabaseVersion)
}

func TestOpenReturnsErrorWhenErrorLoadingMetadata(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")
	require.NoError(t, os.WriteFile(meta_encoder.Encoder{}.MetadataPath(PATH), []byte{}, 0666))

	//-- act
	_, res := vault.Open(PATH)

	//-- assert
	assert.ErrorContains(t, res, "error loading metadata")
}

func TestOpenReturnsVaultAndNoErrorWhenDatabaseIsV3AndLoadsMetadata(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")
	INDEX := index.IndexMap{
		"name": uuid.New(),
	}
	require.NoError(t, database.SaveIndex(db_v3.Database{}, PATH, INDEX))

	VAULT := vault.Vault{
		Path:        PATH,
		MetaEncoder: meta_encoder.Encoder{},
		Meta: meta.Metadata{
			DatabaseVersion: db_v3.VERSION,
		},
	}
	require.NoError(t, VAULT.SaveMetadata())

	//-- act
	res, err := vault.Open(PATH)

	//-- assert
	require.NoError(t, err)
	assert.IsType(t, db_v3.Database{}, res.Database)

	assert.Equal(t, PATH, res.Path)
	assert.Equal(t, INDEX, res.Map)
	assert.Equal(t, VAULT.Meta, res.Meta)
}

func TestOpenReturnsErrorWhenDatabaseVersionUnsupported(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")

	VAULT := vault.Vault{
		Path:        PATH,
		MetaEncoder: meta_encoder.Encoder{},
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
