package output_flow

import (
	"path/filepath"
	"pvault/app/config"
	record_v2 "pvault/vault/record/record/v2"

	"github.com/binarysoupdev/go-commando/logger"

	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

func SaveRecord(cfg config.Config, r record_v2.Record) error {
	err := cfg.ValidateOutputPath()
	if err != nil {
		return errors.Chain(err, "error validating \"config.output_path\"")
	}

	path := filepath.Join(cfg.OutputPath, r.ID.String()+".json")

	err = json.MarshalFilePretty(r, path, "    ")
	if err != nil {
		logger.LogError(err)
		return errors.New("error creating output record")
	}

	logger.Logf("[+] created %s", path)
	style.Create.Printf("[+] %s\n", path)

	return nil
}
