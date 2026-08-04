package config_test

import (
	"path/filepath"
	"pvault/app/config"
	"pvault/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateBackupPathWherePathIsAFileReturnsError(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		BackupPath: filepath.Join(t.TempDir(), "backup.txt"),
	}
	require.NoError(t, util.CreateEmptyFile(cfg.BackupPath))

	//-- act
	res := cfg.ValidateBackupPath()

	//-- assert
	require.ErrorContains(t, res, "path not a directory")
}

func TestValidateBackupPathWhereDirectoryDoesNotExistValid(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		BackupPath: filepath.Join(t.TempDir(), "not_exist"),
	}

	//-- act
	res := cfg.ValidateBackupPath()

	//-- assert
	require.NoError(t, res)
}

func TestValidateBackupPathWhereDirectoryExistsValid(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		BackupPath: t.TempDir(),
	}

	//-- act
	res := cfg.ValidateBackupPath()

	//-- assert
	require.NoError(t, res)
}

func TestValidateOutputPathNotFound(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		OutputPath: filepath.Join(t.TempDir(), "invalid"),
	}

	//-- act
	res := cfg.ValidateOutputPath()

	//-- assert
	require.ErrorContains(t, res, "path not found")
}

func TestValidateOutputPathNotDirectory(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		OutputPath: filepath.Join(t.TempDir(), "file.txt"),
	}
	require.NoError(t, util.CreateEmptyFile(cfg.OutputPath))

	//-- act
	res := cfg.ValidateOutputPath()

	//-- assert
	require.ErrorContains(t, res, "path not a directory")
}

func TestValidateOutputPathValid(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		OutputPath: t.TempDir(),
	}

	//-- act
	res := cfg.ValidateOutputPath()

	//-- assert
	require.NoError(t, res)
}
