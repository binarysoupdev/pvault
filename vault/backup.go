package vault

import (
	"os"

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

	return v.Database.Backup(path, v.Index)
}
