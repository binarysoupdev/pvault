package v1

import (
	"os"
	v2 "pvault/vault/database/version2"
	"pvault/vault/index"
	v1 "pvault/vault/record/version1"

	"github.com/binarysoupdev/go-commando/errors"
)

func (idx Index) Upgrade(m index.IndexMap) error {
	target := v2.NewIndex(idx.Path)

	for name, id := range m {
		legacyFile := idx.RecordPath(id)

		raw, err := os.ReadFile(legacyFile)
		if err != nil {
			continue
		}

		file, err := os.Create(target.RecordPath(id))
		if err != nil {
			return errors.Chain(err, "error creating converted record file")
		}
		defer file.Close()

		v1.EncodeFromLegacy(file, name, raw)
		_ = os.Remove(legacyFile)
	}

	_ = os.Remove(idx.Filepath())
	return nil
}
