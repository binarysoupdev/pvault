package flow_test

import (
	"fmt"
	"path/filepath"
	"pvault/cmd/flow"
	"pvault/vault"
	"pvault/vault/data/version1"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/stretchr/testify/require"
)

func TestOpenVaultVaultOutOfDateReturnsError(t *testing.T) {
	//-- arrange
	PATH := file.CreateEmpty(t, version1.INDEX_FILE)
	const LEGACY_VERSION = 1

	//-- act
	_, res := flow.OpenVault(filepath.Dir(PATH))

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("vault (@v%d) out-of-date", LEGACY_VERSION))
}

func TestOpenVaultInvalidPathReturnsError(t *testing.T) {
	//-- act
	_, res := flow.OpenVault("invalid")

	//-- assert
	require.ErrorContains(t, res, "error opening vault")
}

func TestOpenVaultReturnsVault(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "vault")

	_, err := vault.InitializeNew(PATH)
	require.NoError(t, err)

	//-- act
	_, res := flow.OpenVault(PATH)

	//-- assert
	require.NoError(t, res)
}
