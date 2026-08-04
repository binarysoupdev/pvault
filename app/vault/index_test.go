package vault_test

import (
	"errors"
	"os"
	"pvault/app/vault"
	"pvault/app/vault/database"
	"pvault/app/vault/index"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVaultSaveIndexReturnsErrorWhenDatabaseSaveIndexReturnsError(t *testing.T) {
	//-- arrange
	v := vault.Vault{
		Path: t.TempDir(),
		Database: &database.Mock{
			EncodeIndexError: errors.New(""),
		},
	}

	//-- act
	res := v.SaveIndex()

	//-- assert
	assert.ErrorContains(t, res, "error saving index map")
}

func TestVaultSaveIndexReturnsNoErrorAndSavesIndexMap(t *testing.T) {
	//-- arrange
	mock := &database.Mock{}

	v := vault.Vault{
		Path:     t.TempDir(),
		Database: mock,
		Map: index.IndexMap{
			"name": uuid.New(),
		},
	}

	//-- act
	res := v.SaveIndex()

	//-- assert
	require.NoError(t, res)
	assert.FileExists(t, v.IndexPath())
	assert.Equal(t, v.Map, mock.Index)
}

func TestVaultLoadIndexReturnsErrorWhenDatabaseLoadIndexReturnsError(t *testing.T) {
	//-- arrange
	v := vault.Vault{
		Path: t.TempDir(),
		Database: &database.Mock{
			DecodeIndexError: errors.New(""),
		},
	}

	//-- act
	res := v.LoadIndex()

	//-- assert
	assert.ErrorContains(t, res, "error loading index map")
}

func TestVaultLoadIndexReturnsNoErrorAndLoadsIndexMap(t *testing.T) {
	//-- arrange
	mock := &database.Mock{
		Index: index.IndexMap{
			"name": uuid.New(),
		},
	}

	v := vault.Vault{
		Path:     t.TempDir(),
		Database: mock,
	}
	require.NoError(t, os.WriteFile(v.IndexPath(), []byte{}, 0666))

	//-- act
	res := v.LoadIndex()

	//-- assert
	require.NoError(t, res)
	assert.Equal(t, mock.Index, v.Map)
}
