package v1

import (
	"os"
	index_v2 "pvault/app/vault/index/version2"
	record_v1 "pvault/app/vault/record/version1"

	"github.com/binarysoupdev/go-commando/errors"
)

func (idx Index) Upgrade() (index_v2.Index, error) {
	target := index_v2.NewIndex(idx.Path)

	m, err := idx.LoadMap()
	if err != nil {
		return target, errors.Chain(err, "error loading index")
	}

	err = target.SaveMap(m)
	if err != nil {
		return target, errors.Chain(err, "error saving new index file")
	}

	err = os.Remove(idx.Filepath())
	if err != nil {
		return target, errors.Chain(err, "error removing old index file")
	}

	for name, id := range m {
		legacyFile := idx.RecordPath(id)

		raw, err := os.ReadFile(legacyFile)
		if err != nil {
			continue
		}

		file, err := os.Create(target.RecordPath(id))
		if err != nil {
			continue
		}
		defer file.Close()

		record_v1.EncodeFromLegacy(file, name, raw)
		_ = os.Remove(legacyFile)
	}

	return target, nil
}
