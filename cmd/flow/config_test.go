package flow_test

import (
	"fmt"
	"pvault/cmd/flow"
	"pvault/config"
	"testing"

	cfg "github.com/binarysoupdev/go-commando/config"
	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/tinsel/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigReturnsErrorWhenConfigNotFound(t *testing.T) {
	//-- arrange
	loader := cfg.NewLoader[config.Config]("invalid")

	//-- act
	res := flow.LoadConfig(&loader)

	//-- assert
	require.ErrorContains(t, res, "invalid config path")
}

func TestLoadConfigReturnsErrorWhenUsingInvalidVersion(t *testing.T) {
	//-- arrange
	loader := cfg.NewLoader[config.Config](file.NewPath(t, "config.json"))

	CONFIG := config.Config{
		Version: config.VERSION + 1,
	}
	err := json.MarshalFile(CONFIG, loader.ConfigPath)
	require.NoError(t, err)

	//-- act
	res := flow.LoadConfig(&loader)

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("unsupported version \"%d\"", CONFIG.Version))
}

func TestLoadConfigReturnsNoErrorAndLoadsConfigWhenValid(t *testing.T) {
	//-- arrange
	loader := cfg.NewLoader[config.Config](file.NewPath(t, "config.json"))

	CONFIG := config.Config{
		Version:    config.VERSION,
		BackupPath: "backup/path",
		VaultPath:  "vault/path",
		OutputPath: "output/path",
	}
	err := json.MarshalFile(CONFIG, loader.ConfigPath)
	require.NoError(t, err)

	//-- act
	res := flow.LoadConfig(&loader)

	//-- assert
	require.NoError(t, res)
	assert.Equal(t, CONFIG, loader.Config)
}
