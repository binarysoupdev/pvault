package config_test

import (
	"pvault/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsUnsupportedReturnsTrueWhenVersionLessThanMin(t *testing.T) {
	//-- act
	res := config.IsUnsupported(config.MIN_VERSION - 1)

	//-- assert
	require.True(t, res)
}

func TestIsUnsupportedReturnsTrueWhenVersionGreaterThanCurrent(t *testing.T) {
	//-- act
	res := config.IsUnsupported(config.VERSION + 1)

	//-- assert
	require.True(t, res)
}

func TestIsUnsupportedReturnsFalseWhenVersionValid(t *testing.T) {
	//-- act
	res := config.IsUnsupported(config.VERSION)

	//-- assert
	require.False(t, res)
}

func TestIsOutOfDateReturnsTrueWhenVersionLessThanCurrent(t *testing.T) {
	//-- act
	res := config.IsOutOfDate(config.VERSION - 1)

	//-- assert
	require.True(t, res)
}

func TestIsOutOfDateReturnsFalseWhenVersionEqualsCurrent(t *testing.T) {
	//-- act
	res := config.IsOutOfDate(config.VERSION)

	//-- assert
	require.False(t, res)
}
