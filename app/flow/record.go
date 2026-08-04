package flow

import (
	"path/filepath"
	"pvault/app/config"
	v2 "pvault/app/vault/record/record/v2"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/got-style/style"
)

func SaveOutputRecord(cfg config.Config, r v2.Record) error {
	err := cfg.ValidateOutputPath()
	if err != nil {
		return errors.Chain(err, "error validating output path")
	}

	path := filepath.Join(cfg.OutputPath, r.GetID().String()+".json")

	err = json.MarshalFilePretty(r, path, "    ")
	if err != nil {
		return errors.Chain(err, "error creating output record")
	}

	style.Create.Printf("[+] %s\n", path)
	return nil
}
