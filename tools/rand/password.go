package rand

import (
	"crypto/rand"
	"math/big"
)

func Password(length int) string {
	if length < 1 {
		return ""
	}

	const START = 32 // SPACE
	const END = 126  // ~

	bytes := make([]byte, length)
	max := big.NewInt(END - START)

	for i := range bytes {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err) // rand.Reader should never fail
		}

		bytes[i] = byte(num.Int64() + START)
	}

	return string(bytes)
}
