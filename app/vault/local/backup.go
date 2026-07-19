package local

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

	err = v.backupFile(path, v.Index.Filepath())
	if err != nil {
		return errors.Chain(err, "error backing index file")
	}

	for _, id := range v.Map {
		err := v.backupFile(path, v.Index.RecordPath(id))
		if err != nil {
			continue
		}
	}

	return nil
}

func (Vault) backupFile(dir string, src string) error {
	dest := filepath.Join(dir, filepath.Base(src))

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	if err != nil {
		return err
	}

	return nil
}
