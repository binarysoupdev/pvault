package record_test

import (
	"encoding/json"
	"pvault/app/vault/record"
	"testing"

	"github.com/binarysoupdev/cryptool/crypt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptReturnsPlaintextAndNoErrorWhenPasswordEmpty(t *testing.T) {
	//-- arrange
	RECORD := record.NewMock("name")
	PLAINTEXT, _ := json.Marshal(RECORD)

	//-- act
	res, err := record.Encrypt(RECORD, "")

	//-- arrange
	require.NoError(t, err)
	assert.Equal(t, PLAINTEXT, res)
}

func TestDecryptReturnsRecordAndNoErrorWhenPasswordEmpty(t *testing.T) {
	//-- arrange
	RECORD := record.NewMock("name")
	CIPHERTEXT, _ := json.Marshal(RECORD)

	//-- act
	res, err := record.Decrypt[record.Mock](CIPHERTEXT, "")

	//-- arrange
	require.NoError(t, err)
	assert.Equal(t, RECORD, res)
}

func TestDecryptReturnsErrorWhenSaltTooShort(t *testing.T) {
	//-- arrange
	CIPHERTEXT := []byte{}
	const PASSWORD = "Password123!"

	//-- act
	_, res := record.Decrypt[record.Mock](CIPHERTEXT, PASSWORD)

	//-- arrange
	assert.ErrorContains(t, res, "length too short: 0")
}

func TestDecryptReturnsErrorWhenCiphertextIsOnlySalt(t *testing.T) {
	//-- arrange
	CIPHERTEXT := make([]byte, crypt.SALT_SIZE)
	const PASSWORD = "Password123!"

	//-- act
	_, res := record.Decrypt[record.Mock](CIPHERTEXT, PASSWORD)

	//-- arrange
	assert.ErrorContains(t, res, "length too short: 16")
}

func TestDecryptReturnsErrorWhenErrorUnmarshalingJSON(t *testing.T) {
	//-- arrange
	const PASSWORD = "Password123!"

	c, salt := crypt.NewFromPassword(PASSWORD)
	CIPHERTEXT := c.Encrypt([]byte("foobar"))

	//-- act
	_, res := record.Decrypt[record.Mock](append(salt, CIPHERTEXT...), PASSWORD)

	//-- arrange
	assert.ErrorContains(t, res, "error unmarshaling json")
}

func TestEncryptDecryptReturnsErrorWhenPasswordIncorrect(t *testing.T) {
	//-- arrange
	RECORD := record.NewMock("name")
	const PASSWORD = "Password123!"

	ciphertext, err := record.Encrypt(RECORD, PASSWORD)
	require.NoError(t, err)

	//-- act
	_, res := record.Decrypt[record.Mock](ciphertext, PASSWORD+"x")

	//-- arrange
	assert.ErrorContains(t, res, "message authentication failed")
}

func TestEncryptDecryptReturnsRecordAndNoError(t *testing.T) {
	//-- arrange
	RECORD := record.NewMock("name")
	const PASSWORD = "Password123!"

	ciphertext, err := record.Encrypt(RECORD, PASSWORD)
	require.NoError(t, err)

	//-- act
	res, err := record.Decrypt[record.Mock](ciphertext, PASSWORD)

	//-- arrange
	require.NoError(t, err)
	assert.Equal(t, RECORD, res)
}
