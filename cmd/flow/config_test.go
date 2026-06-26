package flow_test

import (
	"fmt"
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/data"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigConfigNotFoundReturnsError(t *testing.T) {
	//-- arrange
	loader := config.NewLoader[config.Config]("invalid")

	//-- act
	res := flow.LoadConfig(&loader)

	//-- assert
	require.ErrorContains(t, res, "error loading config")
}

func TestLoadConfigInvalidVersionReturnsError(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	loader := config.NewLoader[config.Config](file.NewPath(t, rand.ASCII(10)))

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

func TestLoadConfigOlderVersionReturnsError(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	loader := config.NewLoader[config.Config](file.NewPath(t, rand.ASCII(10)))

	CONFIG := config.Config{
		Version: config.VERSION - 1,
	}
	err := data.SaveJSON(CONFIG, loader.ConfigPath)
	require.NoError(t, err)

	//-- act
	res := flow.LoadConfig(&loader)

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("config version \"%d\" out-of-date", CONFIG.Version))
}

func TestLoadConfigValidLoadsConfig(t *testing.T) {
	//-- arrange
	rand := rand.New(0)
	loader := config.NewLoader[config.Config](file.NewPath(t, rand.ASCII(10)))

	CONFIG := config.Config{
		Version:    config.VERSION,
		VaultPath:  rand.ASCII(15),
		OutputPath: rand.ASCII(15),
	}
	err := data.SaveJSON(CONFIG, loader.ConfigPath)
	require.NoError(t, err)

	//-- act
	res := flow.LoadConfig(&loader)

	//-- assert
	require.NoError(t, res)
	assert.Equal(t, CONFIG, loader.Config)
}
