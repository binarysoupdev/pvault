package flow_test

import (
	"fmt"
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/data"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigReturnsErrorWhenConfigNotFound(t *testing.T) {
	//-- arrange
	loader := config.NewLoader[config.Config]("invalid")

	//-- act
	res := flow.LoadConfig(&loader)

	//-- assert
	require.ErrorContains(t, res, "error loading config")
}

func TestLoadConfigReturnsErrorWhenUsingInvalidVersion(t *testing.T) {
	//-- arrange
	loader := config.NewLoader[config.Config](file.NewPath(t, "config.json"))

	CONFIG := config.Config{
		Version: config.VERSION + 1,
	}
	err := data.SaveJSON(CONFIG, loader.ConfigPath)
	require.NoError(t, err)

	//-- act
	res := flow.LoadConfig(&loader)

	//-- assert
	require.ErrorContains(t, res, "error validating config version")
}

func TestLoadConfigReturnsErrorWhenUsingOlderVersion(t *testing.T) {
	//-- arrange
	loader := config.NewLoader[config.Config](file.NewPath(t, "config.json"))

	CONFIG := config.Config{
		Version: config.VERSION - 1,
	}
	err := data.SaveJSON(CONFIG, loader.ConfigPath)
	require.NoError(t, err)

	//-- act
	res := flow.LoadConfig(&loader)

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("config version [%d] out-of-date", CONFIG.Version))
}

func TestLoadConfigReturnsNoErrorAndLoadsConfigWhenValid(t *testing.T) {
	//-- arrange
	loader := config.NewLoader[config.Config](file.NewPath(t, "config.json"))

	CONFIG := config.Config{
		Version:    config.VERSION,
		VaultPath:  "vault/path",
		OutputPath: "output/path",
	}
	err := data.SaveJSON(CONFIG, loader.ConfigPath)
	require.NoError(t, err)

	//-- act
	res := flow.LoadConfig(&loader)

	//-- assert
	require.NoError(t, res)
	assert.Equal(t, CONFIG, loader.Config)
}
