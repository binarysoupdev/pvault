package local_test

import (
	"pvault/app/vault/index"
	v2 "pvault/app/vault/index/version2"
	"pvault/app/vault/local"
	"testing"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpgradeReturnsErrorWhenVaultIsNotOutOfDate(t *testing.T) {
	//-- arrange
	v := local.Vault{
		Index: &index.Mock{
			Version: local.CURRENT_VERSION,
		},
	}

	//-- act
	res := v.Upgrade()

	//-- arrange
	require.ErrorContains(t, res, "vault is up-to-date")
}

func TestUpgradeReturnsErrorWhenIndexUpgradeReturnsError(t *testing.T) {
	//-- arrange
	v := local.Vault{
		Index: &index.Mock{
			Version:      local.CURRENT_VERSION - 1,
			UpgradeError: errors.New(""),
		},
	}

	//-- act
	res := v.Upgrade()

	//-- arrange
	require.ErrorContains(t, res, "error upgrading index")
}

func TestUpgradeReturnsNoErrorAndSetsNewIndex(t *testing.T) {
	//-- arrange
	INDEX := v2.NewIndex("path")
	v := local.Vault{
		Index: &index.Mock{
			Version:      local.CURRENT_VERSION - 1,
			UpgradeIndex: INDEX,
		},
	}

	//-- act
	res := v.Upgrade()

	//-- arrange
	require.NoError(t, res)
	assert.Equal(t, INDEX, v.Index)
}
