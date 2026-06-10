package config_test

import (
	"pvault/config"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatePathsInvalid(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		VaultPath:  file.NewPath(t, "invalid"),
		OutputPath: file.NewPath(t, "invalid"),
	}

	//-- act
	res := cfg.Validate()

	//-- assert
	require.Error(t, res)

	assert.ErrorContains(t, res, "\"vault_path\" invalid path")
	assert.ErrorContains(t, res, "\"output_path\" invalid path")
}

func TestValidatePathsNoDirs(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		VaultPath:  file.CreateEmpty(t, "file.txt"),
		OutputPath: file.CreateEmpty(t, "file.txt"),
	}

	//-- act
	res := cfg.Validate()

	//-- assert
	require.Error(t, res)

	assert.ErrorContains(t, res, "\"vault_path\" not a directory")
	assert.ErrorContains(t, res, "\"output_path\" not a directory")
}

func TestValidateValid(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		VaultPath:  file.NewPath(t, ""),
		OutputPath: file.NewPath(t, ""),
	}

	//-- act
	res := cfg.Validate()

	//-- assert
	require.NoError(t, res)
}
