package config_test

import (
	"pvault/config"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/stretchr/testify/require"
)

func TestValidateBackupPathWherePathIsAFileReturnsError(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		BackupPath: file.CreateEmpty(t, "backup.txt"),
	}

	//-- act
	res := cfg.ValidateBackupPath()

	//-- assert
	require.ErrorContains(t, res, "path not a directory")
}

func TestValidateBackupPathWhereDirectoryDoesNotExistValid(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		BackupPath: file.NewPath(t, "not_exist"),
	}

	//-- act
	res := cfg.ValidateBackupPath()

	//-- assert
	require.NoError(t, res)
}

func TestValidateBackupPathWhereDirectoryExistsValid(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		BackupPath: file.NewPath(t, ""),
	}

	//-- act
	res := cfg.ValidateBackupPath()

	//-- assert
	require.NoError(t, res)
}

func TestValidateOutputPathNotFound(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		OutputPath: file.NewPath(t, "invalid"),
	}

	//-- act
	res := cfg.ValidateOutputPath()

	//-- assert
	require.ErrorContains(t, res, "path not found")
}

func TestValidateOutputPathNotDirectory(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		OutputPath: file.CreateEmpty(t, "file.txt"),
	}

	//-- act
	res := cfg.ValidateOutputPath()

	//-- assert
	require.ErrorContains(t, res, "path not a directory")
}

func TestValidateOutputPathValid(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		OutputPath: file.NewPath(t, ""),
	}

	//-- act
	res := cfg.ValidateOutputPath()

	//-- assert
	require.NoError(t, res)
}
