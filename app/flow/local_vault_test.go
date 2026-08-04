package flow_test

import (
	"fmt"
	"pvault/app/config"
	"pvault/app/flow"
	"pvault/app/vault/index"
	index_v1 "pvault/app/vault/index/version1"
	"pvault/app/vault/local"
	"regexp"

	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadLocalVaultReturnsErrorWithInvalidPath(t *testing.T) {
	//-- act
	_, res := flow.OpenLocalVault("invalid")

	//-- assert
	require.ErrorContains(t, res, "error opening vault")
}

func TestLoadLocalVaultReturnsErrorWhenVaultOutOfDate(t *testing.T) {
	//-- arrange
	PATH := file.NewPath(t, "")

	err := index_v1.NewIndex(PATH).SaveMap(index.IndexMap{})
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

func TestBackupVaultReturnsErrorWhenBackupPathInvalid(t *testing.T) {
	//-- arrange
	CONFIG := config.Config{
		BackupPath: file.CreateEmpty(t, "invalid.txt"),
	}

	//-- act
	res := flow.BackupVault(local.Vault{}, CONFIG)

	//-- assert
	require.ErrorContains(t, res, "error validating backup path")
}

func TestBackupVaultReturnsNoErrorAndBacksUpVault(t *testing.T) {
	//-- arrange
	DIR_REGEX := regexp.MustCompile(`"([^"]*)"`)

	CONFIG := config.Config{
		VaultPath:  file.NewPath(t, "vault"),
		BackupPath: file.NewPath(t, ""),
	}

	v, err := local.CreateNewVault(CONFIG.VaultPath)
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
