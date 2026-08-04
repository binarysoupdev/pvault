package vault

import (
	"os"
	"path/filepath"
	"pvault/util"

	"github.com/binarysoupdev/go-commando/errors"
)

func (v Vault) Backup(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return errors.Chain(err, "error reading backup directory")
	}

	if !stat.IsDir() {
		return errors.Format("\"%s\" is not a directory", path)
	}

	err = v.backupFile(path, v.MetadataPath())
	if err != nil {
		return errors.Chain(err, "error backing metadata file")
	}

	err = v.backupFile(path, v.DatabaseEncoder.IndexPath(v.Path))
	if err != nil {
		return errors.Chain(err, "error backing index file")
	}
	errs := errors.Errors{}

	for _, id := range v.Map {
		err := v.backupFile(path, v.DatabaseEncoder.RecordPath(v.Path, id))
		if err != nil {
			errs.Add(errors.Chain(err, "error backing record "+id.String()))
			continue
		}
	}

	return errs.Collapse("\n")
}

func (Vault) backupFile(dir string, src string) error {
	return util.CopyFile(filepath.Join(dir, filepath.Base(src)), src)
}
