package vault

import (
	"io"
	"os"
	"path/filepath"

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

	// TODO: move backup to database
	err = v.backupFile(path, v.Database.IndexPath())
	if err != nil {
		return errors.Chain(err, "error backing index file")
	}

	for _, id := range v.Index {
		err := v.backupFile(path, v.Database.RecordPath(id))
		if err != nil {
			return errors.Chain(err, "error backing record")
		}
	}

	return nil
}

func (Vault) backupFile(dir string, src string) error {
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
