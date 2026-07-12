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

func TestConfigCommandLoadConfigReturnsErrorWhenConfigNotFound(t *testing.T) {
	//-- arrange
	cmd := flow.NewConfigCommand(
		json.NewLoader[config.Config]("invalid"),
	)

	//-- act
	res := cmd.LoadConfig()

	//-- assert
	require.ErrorContains(t, res, "invalid config path")
}

func TestConfigCommandLoadConfigReturnsErrorWhenConfigIsInvalidJson(t *testing.T) {
	//-- arrange
	cmd := flow.NewConfigCommand(
		json.NewLoader[config.Config](file.CreateEmpty(t, "invalid.json")),
	)

	//-- act
	res := cmd.LoadConfig()

	//-- assert
	require.ErrorContains(t, res, "error loading config")
}

func TestConfigCommandLoadConfigReturnsErrorWhenUsingInvalidVersion(t *testing.T) {
	//-- arrange
	cmd := flow.NewConfigCommand(
		json.NewLoader[config.Config](file.NewPath(t, "config.json")),
	)

	CONFIG := config.Config{
		Version: config.VERSION + 1,
	}
	err := json.MarshalFile(CONFIG, cmd.ConfigLoader.Path)
	require.NoError(t, err)

	//-- act
	res := cmd.LoadConfig()

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("unsupported config version \"%d\"", CONFIG.Version))
}

func TestConfigCommandLoadConfigReturnsNoErrorAndLoadsConfigWhenValid(t *testing.T) {
	//-- arrange
	cmd := flow.NewConfigCommand(
		json.NewLoader[config.Config](file.NewPath(t, "config.json")),
	)

	CONFIG := config.Config{
		Version:    config.VERSION,
		BackupPath: "backup/path",
		VaultPath:  "vault/path",
		OutputPath: "output/path",
	}
	err := json.MarshalFile(CONFIG, cmd.ConfigLoader.Path)
	require.NoError(t, err)

	//-- act
	res := cmd.LoadConfig()

	//-- assert
	require.NoError(t, res)
	assert.Equal(t, CONFIG, cmd.Config)
}
