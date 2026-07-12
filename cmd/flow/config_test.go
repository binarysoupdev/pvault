package flow_test

import (
	"fmt"
	"pvault/cmd/flow"
	"pvault/config"
	"testing"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigReturnsErrorWhenConfigNotFound(t *testing.T) {
	//-- arrange
	loader := json.NewLoader[config.Config]("invalid")

	//-- act
	_, res := flow.LoadConfig(loader)

	//-- assert
	require.ErrorContains(t, res, "invalid config path")
}

func TestLoadConfigReturnsErrorWhenConfigIsInvalidJson(t *testing.T) {
	//-- arrange
	loader := json.NewLoader[config.Config](file.CreateEmpty(t, "invalid.json"))

	//-- act
	_, res := flow.LoadConfig(loader)

	//-- assert
	require.ErrorContains(t, res, "error loading config")
}

func TestLoadConfigReturnsErrorWhenUsingInvalidVersion(t *testing.T) {
	//-- arrange
	loader := json.NewLoader[config.Config](file.NewPath(t, "config.json"))

	CONFIG := config.Config{
		Version: config.VERSION + 1,
	}
	err := json.MarshalFile(CONFIG, loader.Path)
	require.NoError(t, err)

	//-- act
	_, res := flow.LoadConfig(loader)

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("unsupported version \"%d\"", CONFIG.Version))
}

func TestLoadConfigReturnsConfigAndNoErrorWhenValid(t *testing.T) {
	//-- arrange
	loader := json.NewLoader[config.Config](file.NewPath(t, "config.json"))

	CONFIG := config.Config{
		Version:    config.VERSION,
		BackupPath: "backup/path",
		VaultPath:  "vault/path",
		OutputPath: "output/path",
	}
	err := json.MarshalFile(CONFIG, loader.Path)
	require.NoError(t, err)

	//-- act
	res, err := flow.LoadConfig(loader)

	//-- assert
	require.NoError(t, err)
	assert.Equal(t, CONFIG, res)
}
