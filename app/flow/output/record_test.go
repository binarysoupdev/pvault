package output_flow_test

import (
	"path/filepath"
	"pvault/app/config"
	flow "pvault/app/flow/output"
	record_v2 "pvault/vault/record/record/v2"
	"testing"

	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveRecordReturnsErrorWhenOutputPathInvalid(t *testing.T) {
	//-- arrange
	CONFIG := config.Config{
		OutputPath: "invalid",
	}
	RECORD := record_v2.Record{}

	//-- act
	res := flow.SaveRecord(CONFIG, RECORD)

	//-- assert
	require.ErrorContains(t, res, "error validating \"config.output_path\"")
}

func TestSaveRecordReturnsNoErrorAndSavesJson(t *testing.T) {
	//-- arrange
	CONFIG := config.Config{
		OutputPath: t.TempDir(),
	}
	RECORD := record_v2.NewEmptyRecord("name")

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	res := flow.SaveRecord(CONFIG, RECORD)

	//-- assert
	require.NoError(t, res)
	assert.Contains(t, out.ReadLine(), "[+] "+filepath.Join(CONFIG.OutputPath, RECORD.ID.String()+".json"))
}
