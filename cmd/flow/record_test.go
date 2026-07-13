package flow_test

import (
	"errors"
	"path/filepath"
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/vault"
	"pvault/vault/database"
	"pvault/vault/index"
	v2 "pvault/vault/record/version/v2"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveVaultRecordReturnsErrorWithInvalidRecord(t *testing.T) {
	//-- arrange
	VAULT := vault.Vault{}
	RECORD := v2.Record{}

	//-- act
	res := flow.SaveVaultRecord(VAULT, RECORD)

	//-- arrange
	require.ErrorContains(t, res, "error validating record")
}

func TestSaveVaultRecordReturnsErrorWhenVerifyPasswordDoesNotMatch(t *testing.T) {
	//-- arrange
	VAULT := vault.Vault{}
	RECORD := v2.NewFromName("name")

	const PASSWORD = "Password123!"

	io := pipe.OpenStdio(2, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.Queue("PASSWORD: ", PASSWORD+"x")
	io.EndQueue()

	res := flow.SaveVaultRecord(VAULT, RECORD)

	//-- assert
	require.ErrorContains(t, res, "passwords do not match")
	assert.Contains(t, io.ReadLine(), "New PASSWORD")
	assert.Contains(t, io.ReadLine(), "Verify PASSWORD")
}

func TestSaveVaultRecordReturnsErrorWhenDatabaseSaveRecordFails(t *testing.T) {
	//-- arrange
	VAULT := vault.Vault{
		Database: &database.DatabaseMock{
			SaveRecordError: errors.New(""),
		},
	}
	RECORD := v2.NewFromName("name")

	const PASSWORD = "Password123!"

	io := pipe.OpenStdio(2, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.Queue("PASSWORD: ", PASSWORD)
	io.EndQueue()

	res := flow.SaveVaultRecord(VAULT, RECORD)

	//-- assert
	require.ErrorContains(t, res, "error saving vault record")
	assert.Contains(t, io.ReadLine(), "New PASSWORD")
	assert.Contains(t, io.ReadLine(), "Verify PASSWORD")
}

func TestSaveVaultRecordReturnsNoErrorAndSavesRecordWhenValid(t *testing.T) {
	//-- arrange
	mock := database.DatabaseMock{}

	VAULT := vault.Vault{
		Database: &mock,
		Index:    index.IndexMap{},
	}
	RECORD := v2.NewFromName("name")

	const PASSWORD = "Password123!"

	io := pipe.OpenStdio(2, 3, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.Queue("PASSWORD: ", PASSWORD)
	io.EndQueue()

	res := flow.SaveVaultRecord(VAULT, RECORD)

	//-- assert
	require.NoError(t, res)

	assert.Contains(t, io.ReadLine(), "New PASSWORD")
	assert.Contains(t, io.ReadLine(), "Verify PASSWORD")
	assert.Contains(t, io.ReadLine(), "[+] Saved Record")

	assert.Equal(t, RECORD, mock.Record)
	assert.Contains(t, mock.Index, RECORD.Name)
}

func TestLoadVaultRecordReturnsErrorWhenDatabaseLoadRecordFails(t *testing.T) {
	//-- arrange
	const NAME = "name"
	VAULT := vault.Vault{
		Database: &database.DatabaseMock{
			LoadRecordError: errors.New(""),
		},
		Index: index.IndexMap{
			NAME: uuid.Nil,
		},
	}

	const PASSWORD = "Password123!"

	io := pipe.OpenStdio(1, 1, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.EndQueue()

	_, res := flow.LoadVaultRecord(VAULT, NAME)

	//-- assert
	require.ErrorContains(t, res, "error loading vault record")
	assert.Contains(t, io.ReadLine(), "Enter PASSWORD")
}

func TestLoadVaultRecordReturnsRecordAndNoErrorWhenValid(t *testing.T) {
	//-- arrange
	const NAME = "name"
	mock := database.DatabaseMock{
		Record: v2.NewFromName(NAME),
	}

	VAULT := vault.Vault{
		Database: &mock,
		Index: index.IndexMap{
			NAME: uuid.Nil,
		},
	}
	const PASSWORD = "Password123!"

	io := pipe.OpenStdio(1, 2, false)
	defer io.Close()

	//-- act
	io.Queue("PASSWORD: ", PASSWORD)
	io.EndQueue()

	res, err := flow.LoadVaultRecord(VAULT, NAME)

	//-- assert
	require.NoError(t, err)

	assert.Contains(t, io.ReadLine(), "Enter PASSWORD")
	assert.Contains(t, io.ReadLine(), "[=] Loaded Record")

	assert.Equal(t, res, mock.Record)
}

func TestDeleteVaultRecordReturnsErrorWhenVerifyNameDoesNotMatch(t *testing.T) {
	//-- arrange
	VAULT := vault.Vault{}
	const NAME = "name"

	io := pipe.OpenStdio(1, 1, false)
	defer io.Close()

	//-- act
	io.Queue("NAME: ", NAME+"x")
	io.EndQueue()

	res := flow.DeleteVaultRecord(VAULT, NAME)

	//-- assert
	require.ErrorContains(t, res, "names do not match")
	assert.Contains(t, io.ReadLine(), "Confirm NAME")
}

func TestDeleteVaultRecordReturnsErrorWhenDatabaseDeleteRecordFails(t *testing.T) {
	//-- arrange
	const NAME = "name"
	VAULT := vault.Vault{
		Database: &database.DatabaseMock{
			DeleteRecordError: errors.New(""),
		},
		Index: index.IndexMap{
			NAME: uuid.Nil,
		},
	}

	io := pipe.OpenStdio(1, 1, false)
	defer io.Close()

	//-- act
	io.Queue("NAME: ", NAME)
	io.EndQueue()

	res := flow.DeleteVaultRecord(VAULT, NAME)

	//-- assert
	require.ErrorContains(t, res, "error deleting vault record")
	assert.Contains(t, io.ReadLine(), "Confirm NAME")
}

func TestDeleteVaultRecordReturnsNoErrorAndDeletesRecordWhenValid(t *testing.T) {
	//-- arrange
	const NAME = "name"
	mock := database.DatabaseMock{
		Record: v2.NewFromName(NAME),
	}

	VAULT := vault.Vault{
		Database: &mock,
		Index: index.IndexMap{
			NAME: uuid.Nil,
		},
	}

	io := pipe.OpenStdio(1, 2, false)
	defer io.Close()

	//-- act
	io.Queue("NAME: ", NAME)
	io.EndQueue()

	err := flow.DeleteVaultRecord(VAULT, NAME)

	//-- assert
	require.NoError(t, err)

	assert.Contains(t, io.ReadLine(), "Confirm NAME")
	assert.Contains(t, io.ReadLine(), "[-] Deleted Record")

	assert.NotContains(t, mock.Index, NAME)
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

func TestSaveOutputRecordReturnsNoErrorAndSavesJsonWhenValid(t *testing.T) {
	//-- arrange
	CONFIG := config.Config{
		OutputPath: file.NewPath(t, ""),
	}
	RECORD := v2.NewFromName("name")

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	res := flow.SaveOutputRecord(CONFIG, RECORD)

	//-- assert
	require.NoError(t, res)
	assert.Contains(t, out.ReadLine(), "[+] "+filepath.Join(CONFIG.OutputPath, RECORD.ID.String()+".json"))
}
