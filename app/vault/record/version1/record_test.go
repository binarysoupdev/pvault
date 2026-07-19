package v1_test

import (
	"encoding/binary"
	"fmt"
	v1 "pvault/app/vault/record/version1"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalReturnsErrorWhenVersionIncorrect(t *testing.T) {
	//-- arrange
	const VERSION = v1.VERSION + 1

	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, VERSION)

	//-- act
	_, res := v1.Unmarshal(bytes, "", uuid.Nil)

	//-- assert
	require.ErrorContains(t, res, fmt.Sprintf("incorrect version \"%d\"", VERSION))
}
