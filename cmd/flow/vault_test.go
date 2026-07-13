package flow_test

import (
	"fmt"
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/vault"
	"pvault/vault/database/version1"
	"pvault/vault/index"
	"regexp"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenVaultReturnsErrorWhenVaultOutOfDate(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")

	err := version1.New(PATH).Initialize(index.IndexMap{})
	require.NoError(t, err)

	//-- act
	_, res := flow.OpenVault(PATH)

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("vault (@v%d) out-of-date", version1.VERSION))
}

func TestOpenVaultReturnsErrorWithInvalidPath(t *testing.T) {
	//-- act
	_, res := flow.OpenVault("invalid")

	//-- assert
	require.ErrorContains(t, res, "error opening vault")
}

func TestOpenVaultReturnsVaultAndNoErrorWhenValid(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "vault")

	_, err := vault.InitializeNew(PATH)
	require.NoError(t, err)

	//-- act
	_, res := flow.OpenVault(PATH)

	//-- assert
	require.NoError(t, res)
}

func TestBackupVaultReturnsErrorWhenBackupPathInvalid(t *testing.T) {
	//-- arrange
	CONFIG := config.Config{
		BackupPath: file.CreateEmpty(t, "invalid.txt"),
	}

	//-- act
	res := flow.BackupVault(vault.Vault{}, CONFIG)

	//-- assert
	require.ErrorContains(t, res, "error validating backup path")
}

func TestBackupVaultReturnsNoErrorAndBacksUpVaultWhenValid(t *testing.T) {
	//-- arrange
	DIR_REGEX := regexp.MustCompile(`"([^"]*)"`)

	CONFIG := config.Config{
		VaultPath:  file.NewPath(t, "vault"),
		BackupPath: file.NewPath(t, ""),
	}

	v, err := vault.InitializeNew(CONFIG.VaultPath)
	require.NoError(t, err)

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	res := flow.BackupVault(v, CONFIG)

	//-- assert
	require.NoError(t, res)

	line := out.ReadLine()
	require.Contains(t, line, "[+] Created Backup")

	match := DIR_REGEX.FindStringSubmatch(line)
	require.Len(t, match, 2)
	assert.DirExists(t, match[1])
}
