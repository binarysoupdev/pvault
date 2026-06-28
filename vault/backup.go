package vault

import (
	"io"
	"os"
	"path/filepath"
	"pvault/errors"
)

func (v Vault) Backup(name string) error {
	filenames, err := os.ReadDir(v.Path)
	if err != nil {
		return errors.Chain(err, "error reading vault directory")
	}

	path := filepath.Join(v.Path, name)
	err = os.Mkdir(path, 0755)
	if err != nil {
		return errors.Chain(err, "error creating backup directory")
	}

	for _, filename := range filenames {
		err := copyFile(filepath.Join(v.Path, filename.Name()), filepath.Join(path, filename.Name()))
		if err != nil {
			return errors.Chain(err, "error copying file")
		}
	}

	return nil
}

func copyFile(src, dest string) error {
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
