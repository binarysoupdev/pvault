package flow_test

import (
	"path/filepath"
	"pvault/app/config"
	"pvault/app/flow"
	v2 "pvault/app/vault/record/version2"
	"testing"

	"github.com/binarysoupdev/tinsel/file"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
