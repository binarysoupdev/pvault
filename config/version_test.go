package config_test

import (
	"pvault/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsVersionUnsupportedReturnsTrueWhenVersionLessThanMin(t *testing.T) {
	//-- arrange
	const CURRENT_VERSION = 1
	cfg := config.Config{
		Version: 0,
	}

	//-- act
	res := cfg.Version.IsUnsupported(CURRENT_VERSION)

	//-- assert
	require.True(t, res)
}

func TestIsVersionUnsupportedReturnsTrueWhenVersionGreaterThanCurrent(t *testing.T) {
	//-- arrange
	const CURRENT_VERSION = 1
	cfg := config.Config{
		Version: CURRENT_VERSION + 1,
	}

	//-- act
	res := cfg.Version.IsUnsupported(CURRENT_VERSION)

	//-- assert
	require.True(t, res)
}

func TestIsVersionUnsupportedReturnsFalseWhenVersionValid(t *testing.T) {
	//-- arrange
	const CURRENT_VERSION = 1
	cfg := config.Config{
		Version: CURRENT_VERSION,
	}

	//-- act
	res := cfg.Version.IsUnsupported(CURRENT_VERSION)

	//-- assert
	require.False(t, res)
}

func TestIsVersionOutOfDateReturnsTrueWhenVersionLessThanCurrent(t *testing.T) {
	//-- arrange
	const CURRENT_VERSION = 1
	cfg := config.Config{
		Version: CURRENT_VERSION - 1,
	}

	//-- act
	res := cfg.Version.IsOutOfDate(CURRENT_VERSION)

	//-- assert
	require.True(t, res)
}

func TestIsVersionOutOfDateReturnsFalseWhenVersionEqualsCurrent(t *testing.T) {
	//-- arrange
	const CURRENT_VERSION = 1
	cfg := config.Config{
		Version: CURRENT_VERSION,
	}

	//-- act
	res := cfg.Version.IsOutOfDate(CURRENT_VERSION)

	//-- assert
	require.False(t, res)
}
