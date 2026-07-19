package v2_test

import (
	"encoding/binary"
	"fmt"
	v2 "pvault/app/vault/record/version2"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalReturnsErrorWhenVersionIncorrect(t *testing.T) {
	//-- arrange
	const VERSION = v2.VERSION + 1

	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, VERSION)

	//-- act
	_, res := v2.Unmarshal(bytes, "")

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("incorrect version \"%d\"", VERSION))
}
