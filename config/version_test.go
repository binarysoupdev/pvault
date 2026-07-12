package config_test

import (
	"pvault/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsVersionUnsupportedReturnsTrueWhenVersionLessThanMin(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		Version: 0,
	}

	//-- act
	res := cfg.IsVersionUnsupported()

	//-- assert
	require.True(t, res)
}

func TestIsVersionUnsupportedReturnsTrueWhenVersionGreaterThanCurrent(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		Version: config.VERSION + 1,
	}

	//-- act
	res := cfg.IsVersionUnsupported()

	//-- assert
	require.True(t, res)
}

func TestIsVersionUnsupportedReturnsFalseWhenVersionValid(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		Version: config.VERSION,
	}

	//-- act
	res := cfg.IsVersionUnsupported()

	//-- assert
	require.False(t, res)
}

func TestIsVersionOutOfDateReturnsTrueWhenVersionLessThanCurrent(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		Version: config.VERSION - 1,
	}

	//-- act
	res := cfg.IsVersionOutOfDate()

	//-- assert
	require.True(t, res)
}

func TestIsVersionOutOfDateReturnsFalseWhenVersionEqualsCurrent(t *testing.T) {
	//-- arrange
	cfg := config.Config{
		Version: config.VERSION,
	}

	//-- act
	res := cfg.IsVersionOutOfDate()

	//-- assert
	require.False(t, res)
}
