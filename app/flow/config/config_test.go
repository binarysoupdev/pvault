package config_flow_test

import (
	"fmt"
	"path/filepath"
	"pvault/app/config"
	config_flow "pvault/app/flow/config"
	"pvault/util"
	"testing"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigCommandBaseLoadConfigReturnsErrorWhenConfigNotFound(t *testing.T) {
	//-- arrange
	cmd := command.NewConfigCommand(
		json.NewLoader[config.Config]("invalid"),
	)

	//-- act
	res := config_flow.LoadConfig(&cmd)

	//-- assert
	require.ErrorContains(t, res, "invalid config path")
}

func TestConfigCommandBaseLoadConfigReturnsErrorWhenConfigIsInvalidJson(t *testing.T) {
	//-- arrange
	cmd := command.NewConfigCommand(
		json.NewLoader[config.Config](filepath.Join(t.TempDir(), "invalid.json")),
	)
	require.NoError(t, util.CreateEmptyFile(cmd.ConfigLoader.Path))

	//-- act
	res := config_flow.LoadConfig(&cmd)

	//-- assert
	require.ErrorContains(t, res, "error loading config")
}

func TestConfigCommandBaseLoadConfigReturnsErrorWhenUsingInvalidVersion(t *testing.T) {
	//-- arrange
	cmd := command.NewConfigCommand(
		json.NewLoader[config.Config](filepath.Join(t.TempDir(), "config.json")),
	)

	CONFIG := config.Config{
		Version: config.VERSION + 1,
	}
	require.NoError(t, json.MarshalFile(CONFIG, cmd.ConfigLoader.Path))

	//-- act
	res := config_flow.LoadConfig(&cmd)

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("unsupported config version \"%d\"", CONFIG.Version))
}

func TestConfigCommandBaseLoadConfigReturnsNoErrorAndLoadsConfigWhenValid(t *testing.T) {
	//-- arrange
	cmd := command.NewConfigCommand(
		json.NewLoader[config.Config](filepath.Join(t.TempDir(), "config.json")),
	)

	CONFIG := config.Config{
		Version:    config.VERSION,
		BackupPath: "backup/path",
		VaultPath:  "vault/path",
		OutputPath: "output/path",
	}
	require.NoError(t, json.MarshalFile(CONFIG, cmd.ConfigLoader.Path))

	//-- act
	res := config_flow.LoadConfig(&cmd)

	//-- assert
	require.NoError(t, res)
	assert.Equal(t, CONFIG, cmd.Config)
}
