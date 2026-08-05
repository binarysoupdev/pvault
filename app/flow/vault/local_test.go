package vault_test

import (
	"fmt"
	"path/filepath"
	"pvault/app/config"
	flow "pvault/app/flow/vault"
	"pvault/app/vault"
	"pvault/app/vault/database"
	db_v1 "pvault/app/vault/database/encoder/legacy/v1"
	"pvault/app/vault/index"
	"pvault/app/vault/meta"
	meta_v1 "pvault/app/vault/meta/encoder/v1"
	"pvault/util"
	"regexp"

	"testing"

	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenVaultReturnsErrorWithInvalidPath(t *testing.T) {
	//-- act
	_, res := flow.OpenVault("invalid")

	//-- assert
	require.ErrorContains(t, res, "error opening vault")
}

func TestOpenVaultReturnsErrorWhenVaultOutOfDate(t *testing.T) {
	//-- arrange
	PATH := t.TempDir()
	DATABASE := db_v1.Encoder{}

	META := meta.Metadata{
		DatabaseVersion: DATABASE.GetVersion(),
	}
	require.NoError(t, meta.SaveMetadata(meta_v1.Encoder{}, PATH, META))
	require.NoError(t, database.SaveIndex(DATABASE, PATH, index.IndexMap{}))

	//-- act
	_, res := flow.OpenVault(PATH)

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("vault (@v%d) out-of-date", META.DatabaseVersion))
}

func TestOpenVaultReturnsVaultAndNoError(t *testing.T) {
	//-- arrange
	PATH := filepath.Join(t.TempDir(), "vault")

	_, err := vault.InitializeNew(PATH, "")
	require.NoError(t, err)

	//-- act
	_, res := flow.OpenVault(PATH)

	//-- assert
	require.NoError(t, res)
}

func TestBackupVaultReturnsErrorWhenBackupPathInvalid(t *testing.T) {
	//-- arrange
	CONFIG := config.Config{
		BackupPath: filepath.Join(t.TempDir(), "invalid.txt"),
	}
	require.NoError(t, util.CreateEmptyFile(CONFIG.BackupPath))

	//-- act
	res := flow.BackupVault(vault.Vault{}, CONFIG)

	//-- assert
	require.ErrorContains(t, res, "error validating \"config.backup_path\"")
}

func TestBackupVaultReturnsNoErrorAndBacksUpVault(t *testing.T) {
	//-- arrange
	DIR_REGEX := regexp.MustCompile(`"([^"]*)"`)

	CONFIG := config.Config{
		VaultPath:  filepath.Join(t.TempDir(), "vault"),
		BackupPath: t.TempDir(),
	}

	v, err := vault.InitializeNew(CONFIG.VaultPath, "")
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
