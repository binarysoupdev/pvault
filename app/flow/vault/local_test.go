package vault_flow_test

import (
	"fmt"
	"path/filepath"
	"pvault/app/config"
	vault_flow "pvault/app/flow/vault"
	"pvault/util"
	"pvault/vault"
	"pvault/vault/database"
	db_v1 "pvault/vault/database/encoder/legacy/v1"
	db_v2 "pvault/vault/database/encoder/legacy/v2"
	"pvault/vault/index"
	"pvault/vault/meta"
	meta_v1 "pvault/vault/meta/encoder/v1"
	"regexp"

	"testing"

	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenLegacyVaultReturnsErrorWithInvalidPath(t *testing.T) {
	//-- act
	_, res := vault_flow.OpenLegacyVault("invalid")

	//-- assert
	require.ErrorContains(t, res, "vault not found")
}

func TestOpenLegacyVaultReturnsErrorWithInvalidVault(t *testing.T) {
	//-- arrange
	PATH := t.TempDir()
	require.NoError(t, util.CreateEmptyFile(db_v2.Encoder{}.IndexPath(PATH)))

	//-- act
	_, res := vault_flow.OpenLegacyVault(PATH)

	//-- assert
	require.ErrorContains(t, res, "error opening vault")
}

func TestOpenLegacyVaultReturnsVaultAndNoError(t *testing.T) {
	//-- arrange
	PATH := t.TempDir()
	DATABASE := db_v1.Encoder{}

	META := meta.Metadata{
		DatabaseVersion: DATABASE.GetVersion(),
	}
	require.NoError(t, meta.SaveMetadata(meta_v1.Encoder{}, PATH, META))
	require.NoError(t, database.SaveIndex(DATABASE, PATH, index.IndexMap{}))

	//-- act
	res, err := vault_flow.OpenLegacyVault(PATH)

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, DATABASE, res.DatabaseEncoder)
}

func TestOpenCurrentVaultReturnsErrorWithInvalidPath(t *testing.T) {
	//-- act
	_, res := vault_flow.OpenCurrentVault("invalid")

	//-- assert
	require.ErrorContains(t, res, "vault not found")
}

func TestOpenCurrentVaultReturnsErrorWithInvalidVault(t *testing.T) {
	//-- arrange
	PATH := t.TempDir()
	require.NoError(t, util.CreateEmptyFile(db_v2.Encoder{}.IndexPath(PATH)))

	//-- act
	_, res := vault_flow.OpenCurrentVault(PATH)

	//-- assert
	require.ErrorContains(t, res, "error opening vault")
}

func TestOpenCurrentVaultReturnsErrorWhenVaultOutOfDate(t *testing.T) {
	//-- arrange
	PATH := t.TempDir()
	DATABASE := db_v1.Encoder{}

	META := meta.Metadata{
		DatabaseVersion: DATABASE.GetVersion(),
	}
	require.NoError(t, meta.SaveMetadata(meta_v1.Encoder{}, PATH, META))
	require.NoError(t, database.SaveIndex(DATABASE, PATH, index.IndexMap{}))

	//-- act
	_, res := vault_flow.OpenCurrentVault(PATH)

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("vault (@v%d) out-of-date", META.DatabaseVersion))
}

func TestOpenCurrentVaultReturnsVaultAndNoErrorWhenVaultUpToDate(t *testing.T) {
	//-- arrange
	PATH := filepath.Join(t.TempDir(), "vault")

	VAULT, err := vault.InitializeNew(PATH, "")
	require.NoError(t, err)

	//-- act
	res, err := vault_flow.OpenCurrentVault(PATH)

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, VAULT.DatabaseEncoder, res.DatabaseEncoder)
}

func TestBackupVaultReturnsErrorWhenBackupPathInvalid(t *testing.T) {
	//-- arrange
	CONFIG := config.Config{
		BackupPath: filepath.Join(t.TempDir(), "invalid.txt"),
	}
	require.NoError(t, util.CreateEmptyFile(CONFIG.BackupPath))

	//-- act
	res := vault_flow.BackupVault(vault.Vault{}, CONFIG)

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
	res := vault_flow.BackupVault(v, CONFIG)

	//-- assert
	require.NoError(t, res)

	line := out.ReadLine()
	require.Contains(t, line, "[+] Created Backup")

	match := DIR_REGEX.FindStringSubmatch(line)
	require.Len(t, match, 2)
	assert.DirExists(t, match[1])
}
