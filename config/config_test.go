package config_test

import (
	"pvault/config"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/stretchr/testify/require"
)

func TestValidateOutputPathNotFound(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		OutputPath: file.NewPath(t, "invalid"),
	}

	//-- act
	res := cfg.ValidateOutputPath()

	//-- assert
	require.ErrorContains(t, res, "error loading file stats")
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
