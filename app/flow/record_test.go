package flow_test

import (
	"errors"
	"path/filepath"
	"pvault/app/flow"
	"pvault/app/vault"
	v2 "pvault/app/vault/record/version2"
	"pvault/config"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveVaultRecordReturnsErrorWithInvalidRecord(t *testing.T) {
	//-- arrange
	RECORD := v2.NewEmptyRecord("name")
	v := &vault.Mock{
		ValidateRecordError: errors.New(""),
	}

	//-- act
	res := flow.SaveVaultRecord(v, RECORD)

	//-- arrange
	require.ErrorContains(t, res, "error validating record")
	assert.Equal(t, RECORD, v.Record)
}

func TestSaveVaultRecordReturnsErrorWhenVerifyPasswordDoesNotMatch(t *testing.T) {
	//-- arrange
	const PASSWORD = "Password123!"

	io := pipe.OpenStdio(2, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.Queue("PASSWORD: ", PASSWORD+"x")
	io.EndQueue()

	res := flow.SaveVaultRecord(&vault.Mock{}, v2.Record{})

	//-- assert
	require.ErrorContains(t, res, "passwords do not match")
	assert.Contains(t, io.ReadLine(), "New PASSWORD")
	assert.Contains(t, io.ReadLine(), "Verify PASSWORD")
}

func TestSaveVaultRecordReturnsErrorWhenVaultSaveRecordReturnsError(t *testing.T) {
	//-- arrange
	RECORD := v2.NewEmptyRecord("name")
	v := &vault.Mock{
		SaveRecordError: errors.New(""),
	}

	const PASSWORD = "Password123!"

	io := pipe.OpenStdio(2, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.Queue("PASSWORD: ", PASSWORD)
	io.EndQueue()

	res := flow.SaveVaultRecord(v, RECORD)

	//-- assert
	require.ErrorContains(t, res, "error saving vault record")
	assert.Contains(t, io.ReadLine(), "New PASSWORD")
	assert.Contains(t, io.ReadLine(), "Verify PASSWORD")

	assert.Equal(t, RECORD, v.Record)
	assert.Equal(t, PASSWORD, v.PasswordParam)
}

func TestSaveVaultRecordReturnsNoErrorAndSavesRecord(t *testing.T) {
	//-- arrange
	RECORD := v2.NewEmptyRecord("name")
	v := &vault.Mock{}

	const PASSWORD = "Password123!"

	io := pipe.OpenStdio(2, 3, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.Queue("PASSWORD: ", PASSWORD)
	io.EndQueue()

	res := flow.SaveVaultRecord(v, RECORD)

	//-- assert
	require.NoError(t, res)

	assert.Contains(t, io.ReadLine(), "New PASSWORD")
	assert.Contains(t, io.ReadLine(), "Verify PASSWORD")
	assert.Contains(t, io.ReadLine(), "[+] Saved Record")

	assert.Equal(t, RECORD, v.Record)
	assert.Equal(t, PASSWORD, v.PasswordParam)
}

func TestLoadVaultRecordReturnsErrorWhenVaultLoadRecordReturnsError(t *testing.T) {
	//-- arrange
	const NAME = "name"
	v := &vault.Mock{
		LoadRecordError: errors.New(""),
	}

	const PASSWORD = "Password123!"

	io := pipe.OpenStdio(1, 1, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.EndQueue()

	_, res := flow.LoadVaultRecord(v, NAME)

	//-- assert
	require.ErrorContains(t, res, "error loading vault record")
	assert.Contains(t, io.ReadLine(), "Enter PASSWORD")

	assert.Equal(t, NAME, v.NameParam)
	assert.Equal(t, PASSWORD, v.PasswordParam)
}

func TestLoadVaultRecordReturnsRecordAndNoError(t *testing.T) {
	//-- arrange
	const NAME = "name"
	v := &vault.Mock{
		Record: v2.NewEmptyRecord(NAME),
	}

	const PASSWORD = "Password123!"

	io := pipe.OpenStdio(1, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.EndQueue()

	res, err := flow.LoadVaultRecord(v, NAME)

	//-- assert
	require.NoError(t, err)

	assert.Contains(t, io.ReadLine(), "Enter PASSWORD")
	assert.Contains(t, io.ReadLine(), "[=] Loaded Record")

	assert.Equal(t, v.Record, res)
	assert.Equal(t, NAME, v.NameParam)
	assert.Equal(t, PASSWORD, v.PasswordParam)
}

func TestDeleteVaultRecordReturnsErrorWhenVerifyNameDoesNotMatch(t *testing.T) {
	//-- arrange
	const NAME = "name"

	io := pipe.OpenStdio(1, 1, false)
	defer io.Close()

	//-- act
	io.Queue("NAME: ", NAME+"x")
	io.EndQueue()

	res := flow.DeleteVaultRecord(&vault.Mock{}, NAME)

	//-- assert
	require.ErrorContains(t, res, "names do not match")
	assert.Contains(t, io.ReadLine(), "Confirm NAME")
}

func TestDeleteVaultRecordReturnsErrorWhenVaultDeleteRecordReturnsError(t *testing.T) {
	//-- arrange
	const NAME = "name"
	v := &vault.Mock{
		Record:            v2.NewEmptyRecord(NAME),
		DeleteRecordError: errors.New(""),
	}

	io := pipe.OpenStdio(1, 1, false)
	defer io.Close()

	//-- act
	io.Queue("NAME: ", NAME)
	io.EndQueue()

	res := flow.DeleteVaultRecord(v, NAME)

	//-- assert
	require.ErrorContains(t, res, "error deleting vault record")
	assert.Contains(t, io.ReadLine(), "Confirm NAME")

	assert.Equal(t, NAME, v.NameParam)
}

func TestDeleteVaultRecordReturnsNoErrorAndDeletesRecord(t *testing.T) {
	//-- arrange
	const NAME = "name"
	v := &vault.Mock{
		Record: v2.NewEmptyRecord(NAME),
	}

	io := pipe.OpenStdio(1, 2, false)
	defer io.Close()

	//-- act
	io.Queue("NAME: ", NAME)
	io.EndQueue()

	err := flow.DeleteVaultRecord(v, NAME)

	//-- assert
	require.NoError(t, err)

	assert.Contains(t, io.ReadLine(), "Confirm NAME")
	assert.Contains(t, io.ReadLine(), "[-] Deleted Record: "+v.Record.GetID().String())

	assert.Equal(t, NAME, v.NameParam)
}

func TestSaveOutputRecordReturnsErrorWhenOutputPathInvalid(t *testing.T) {
	//-- arrange
	CONFIG := config.Config{
		OutputPath: "invalid",
	}
	RECORD := v2.Record{}

	//-- act
	res := flow.SaveOutputRecord(CONFIG, RECORD)

	//-- assert
	require.ErrorContains(t, res, "error validating output path")
}

func TestSaveOutputRecordReturnsNoErrorAndSavesJson(t *testing.T) {
	//-- arrange
	CONFIG := config.Config{
		OutputPath: file.NewPath(t, ""),
	}
	RECORD := v2.NewEmptyRecord("name")

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	res := flow.SaveOutputRecord(CONFIG, RECORD)

	//-- assert
	require.NoError(t, res)
	assert.Contains(t, out.ReadLine(), "[+] "+filepath.Join(CONFIG.OutputPath, RECORD.ID.String()+".json"))
}
