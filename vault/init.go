package vault

import (
	"os"
	"path/filepath"
	"pvault/errors"
	"pvault/vault/data"
	"pvault/vault/index"
)

const INDEX_FILE = "index.bin"

func (v Vault) NewCurrentDatabase() data.DatabaseV2 {
	return data.NewDatabaseV2(filepath.Join(v.Path, INDEX_FILE))
}

func InitializeNew(path string) (Vault, error) {
	_, err := os.Stat(path)
	if err == nil || !os.IsNotExist(err) {
		return Vault{}, errors.New("vault path already exists")
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return Vault{}, errors.Chain(err, "error creating vault directory")
	}

	v := Vault{
		Path:  path,
		Index: index.IndexMap{},
	}
	v.Database = v.NewCurrentDatabase()

	err = v.Database.SaveIndex(v.Index)
	if err != nil {
		return v, errors.Chain(err, "error saving index file")
	}

	return v, nil
}

func (v *Vault) ReloadIndex() error {
	var err error

	v.Index, err = v.Database.LoadIndex()
	if err != nil {
		return errors.Chain(err, "error loading index")
	}

	return nil
}
