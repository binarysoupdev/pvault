package database

import (
	"io"
	"os"
	"path/filepath"
	"pvault/vault/index"

	"github.com/binarysoupdev/go-commando/errors"
)

func (db Database) Backup(path string, idx index.IndexMap) error {
	stat, err := os.Stat(path)
	if err != nil {
		return errors.Chain(err, "error reading backup directory")
	}

	if !stat.IsDir() {
		return errors.Format("\"%s\" is not a directory", path)
	}

	err = db.backupFile(path, db.Filepath())
	if err != nil {
		return errors.Chain(err, "error backing index file")
	}

	for _, id := range idx {
		err := db.backupFile(path, db.RecordPath(id))
		if err != nil {
			return errors.Chain(err, "error backing record")
		}
	}

	return nil
}

func (Database) backupFile(dir string, src string) error {
	dest := filepath.Join(dir, filepath.Base(src))

	s, err := os.Open(src)
	if err != nil {
		return errors.Chain(err, "error opening source file")
	}
	defer s.Close()

	d, err := os.Create(dest)
	if err != nil {
		return errors.Chain(err, "error creating destination file")
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	if err != nil {
		return errors.Chain(err, "error copying data")
	}

	return nil
}
