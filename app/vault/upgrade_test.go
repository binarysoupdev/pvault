package vault_test

import (
	"pvault/app/vault"
	"pvault/app/vault/database"
	"pvault/app/vault/meta"
	"testing"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVaultUpgradeReturnsErrorWhenVaultIsNotOutOfDate(t *testing.T) {
	//-- arrange
	v := vault.Vault{
		Meta: meta.Metadata{
			DatabaseVersion: database.CURRENT_VERSION,
		},
	}

	//-- act
	res := v.Upgrade()

	//-- arrange
	require.ErrorContains(t, res, "vault is up-to-date")
}

func TestVaultUpgradeReturnsErrorWhenDatabaseUpgradeReturnsError(t *testing.T) {
	//-- arrange
	v := vault.Vault{
		Meta: meta.Metadata{
			DatabaseVersion: database.CURRENT_VERSION - 1,
		},
		DatabaseEncoder: &database.EncoderMock{
			UpgradeError: errors.New(""),
		},
	}

	//-- act
	res := v.Upgrade()

	//-- arrange
	require.ErrorContains(t, res, "error upgrading database")
}

func TestVaultUpgradeReturnsNoErrorAndSetsNewDatabase(t *testing.T) {
	//-- arrange
	mock := &database.EncoderMock{}

	v := vault.Vault{
		Path: t.TempDir(),
		Meta: meta.Metadata{
			DatabaseVersion: database.CURRENT_VERSION - 1,
		},
		DatabaseEncoder: mock,
		MetaEncoder:     &meta.EncoderMock{},
	}

	//-- act
	res := v.Upgrade()

	//-- arrange
	require.NoError(t, res)
	assert.Equal(t, mock.UpgradedEncoder, v.DatabaseEncoder)
}
