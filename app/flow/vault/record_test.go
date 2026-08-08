package vault_flow_test

import (
	"errors"
	vault_flow "pvault/app/flow/vault"
	record_v2 "pvault/app/vault/record/record/v2"
	"testing"

	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveRecordReturnsErrorWithInvalidRecord(t *testing.T) {
	//-- arrange
	RECORD := record_v2.NewEmptyRecord("name")
	v := &vault_flow.VaultMock{
		ValidateRecordError: errors.New(""),
	}

	//-- act
	res := vault_flow.SaveRecord(v, RECORD)

	//-- arrange
	require.ErrorContains(t, res, "error validating record")
	assert.Equal(t, RECORD, v.Record)
}

func TestSaveRecordReturnsErrorWhenVerifyPasswordDoesNotMatch(t *testing.T) {
	//-- arrange
	const PASSWORD = "Password123!"

	io := pipe.OpenStdio(2, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.Queue("PASSWORD: ", PASSWORD+"x")
	io.EndQueue()

	res := vault_flow.SaveRecord(&vault_flow.VaultMock{}, record_v2.Record{})

	//-- assert
	require.ErrorContains(t, res, "passwords do not match")
	assert.Contains(t, io.ReadLine(), "New PASSWORD")
	assert.Contains(t, io.ReadLine(), "Verify PASSWORD")
}

func TestSaveRecordReturnsErrorWhenVaultSaveRecordReturnsError(t *testing.T) {
	//-- arrange
	RECORD := record_v2.NewEmptyRecord("name")
	v := &vault_flow.VaultMock{
		SaveRecordError: errors.New(""),
	}

	const PASSWORD = "Password123!"

	io := pipe.OpenStdio(2, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.Queue("PASSWORD: ", PASSWORD)
	io.EndQueue()

	res := vault_flow.SaveRecord(v, RECORD)

	//-- assert
	require.ErrorContains(t, res, "error saving vault record")
	assert.Contains(t, io.ReadLine(), "New PASSWORD")
	assert.Contains(t, io.ReadLine(), "Verify PASSWORD")

	assert.Equal(t, RECORD, v.Record)
	assert.Equal(t, PASSWORD, v.PasswordParam)
}

func TestSaveRecordReturnsNoErrorAndSavesRecord(t *testing.T) {
	//-- arrange
	RECORD := record_v2.NewEmptyRecord("name")
	v := &vault_flow.VaultMock{}

	const PASSWORD = "Password123!"

	io := pipe.OpenStdio(2, 3, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.Queue("PASSWORD: ", PASSWORD)
	io.EndQueue()

	res := vault_flow.SaveRecord(v, RECORD)

	//-- assert
	require.NoError(t, res)

	assert.Contains(t, io.ReadLine(), "New PASSWORD")
	assert.Contains(t, io.ReadLine(), "Verify PASSWORD")
	assert.Contains(t, io.ReadLine(), "[+] Saved Record")

	assert.Equal(t, RECORD, v.Record)
	assert.Equal(t, PASSWORD, v.PasswordParam)
}

func TestLoadRecordReturnsErrorWhenVaultLoadRecordReturnsError(t *testing.T) {
	//-- arrange
	const NAME = "name"
	v := &vault_flow.VaultMock{
		LoadRecordError: errors.New(""),
	}

	const PASSWORD = "Password123!"

	io := pipe.OpenStdio(1, 1, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.EndQueue()

	_, res := vault_flow.LoadRecord(v, NAME)

	//-- assert
	require.ErrorContains(t, res, "error loading vault record")
	assert.Contains(t, io.ReadLine(), "Enter PASSWORD")

	assert.Equal(t, NAME, v.NameParam)
	assert.Equal(t, PASSWORD, v.PasswordParam)
}

func TestLoadRecordReturnsRecordAndNoError(t *testing.T) {
	//-- arrange
	const NAME = "name"
	v := &vault_flow.VaultMock{
		Record: record_v2.NewEmptyRecord(NAME),
	}

	const PASSWORD = "Password123!"

	io := pipe.OpenStdio(1, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.EndQueue()

	res, err := vault_flow.LoadRecord(v, NAME)

	//-- assert
	require.NoError(t, err)

	assert.Contains(t, io.ReadLine(), "Enter PASSWORD")
	assert.Contains(t, io.ReadLine(), "[=] Loaded Record")

	assert.Equal(t, v.Record, res)
	assert.Equal(t, NAME, v.NameParam)
	assert.Equal(t, PASSWORD, v.PasswordParam)
}

func TestDeleteRecordReturnsErrorWhenVerifyNameDoesNotMatch(t *testing.T) {
	//-- arrange
	const NAME = "name"

	io := pipe.OpenStdio(1, 1, false)
	defer io.Close()

	//-- act
	io.Queue("NAME: ", NAME+"x")
	io.EndQueue()

	res := vault_flow.DeleteRecord(&vault_flow.VaultMock{}, NAME)

	//-- assert
	require.ErrorContains(t, res, "names do not match")
	assert.Contains(t, io.ReadLine(), "Confirm NAME")
}

func TestDeleteRecordReturnsErrorWhenVaultDeleteRecordReturnsError(t *testing.T) {
	//-- arrange
	const NAME = "name"
	v := &vault_flow.VaultMock{
		Record:            record_v2.NewEmptyRecord(NAME),
		DeleteRecordError: errors.New(""),
	}

	io := pipe.OpenStdio(1, 1, false)
	defer io.Close()

	//-- act
	io.Queue("NAME: ", NAME)
	io.EndQueue()

	res := vault_flow.DeleteRecord(v, NAME)

	//-- assert
	require.ErrorContains(t, res, "error deleting vault record")
	assert.Contains(t, io.ReadLine(), "Confirm NAME")

	assert.Equal(t, NAME, v.NameParam)
}

func TestDeleteRecordReturnsNoErrorAndDeletesRecord(t *testing.T) {
	//-- arrange
	const NAME = "name"
	v := &vault_flow.VaultMock{
		Record: record_v2.NewEmptyRecord(NAME),
	}

	io := pipe.OpenStdio(1, 2, false)
	defer io.Close()

	//-- act
	io.Queue("NAME: ", NAME)
	io.EndQueue()

	err := vault_flow.DeleteRecord(v, NAME)

	//-- assert
	require.NoError(t, err)

	assert.Contains(t, io.ReadLine(), "Confirm NAME")
	assert.Contains(t, io.ReadLine(), "[-] Deleted Record: "+v.Record.GetID().String())

	assert.Equal(t, NAME, v.NameParam)
}
