package config_test

import (
	"pvault/config"
	"pvault/data"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigFileNotFound(t *testing.T) {
	//-- arrange
	loader := config.NewLoader[config.Config]("invalid")

	//-- act
	res := loader.LoadConfig()

	//-- assert
	require.ErrorContains(t, res, "error loading config JSON")
}

func TestLoadConfigLoadsConfig(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	loader := config.NewLoader[config.Config](file.NewPath(t, rand.ASCII(10)))

	CONFIG := config.Config{
		Version:    rand.ASCII(3),
		VaultPath:  rand.ASCII(15),
		OutputPath: rand.ASCII(15),
	}

	err := data.SaveJSON(CONFIG, loader.ConfigPath)
	require.NoError(t, err)

	//-- act
	res := loader.LoadConfig()

	//-- assert
	require.NoError(t, res)
	assert.Equal(t, CONFIG, loader.Config)
}
