package flow_test

import (
	"fmt"
	"pvault/app/flow"
	"pvault/app/vault/data"
	index_v1 "pvault/app/vault/index/version1"
	"pvault/app/vault/local"

	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/stretchr/testify/require"
)

func TestLoadLocalVaultReturnsErrorWithInvalidPath(t *testing.T) {
	//-- act
	_, res := flow.OpenLocalVault("invalid")

	//-- assert
	require.ErrorContains(t, res, "error loading vault")
}

func TestLoadLocalVaultReturnsErrorWhenVaultOutOfDate(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")

	err := index_v1.NewIndex(PATH).SaveMap(data.NameMap{})
	require.NoError(t, err)

	//-- act
	_, res := flow.OpenLocalVault(PATH)

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("vault (@v%d) out-of-date", index_v1.VERSION))
}

func TestLoadLocalVaultReturnsVaultAndNoError(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "vault")

	_, err := local.CreateNewVault(PATH)
	require.NoError(t, err)

	//-- act
	_, res := flow.OpenLocalVault(PATH)

	//-- assert
	require.NoError(t, res)
}
