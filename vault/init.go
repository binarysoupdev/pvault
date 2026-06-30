package vault

import (
	"os"
	"pvault/errors"
	"pvault/vault/data/version2"
	"pvault/vault/index"
)

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

	v.Database, err = v.initNewDatabaseVersion2()
	if err != nil {
		return Vault{}, err
	}

	return v, nil
}

func (v Vault) initNewDatabaseVersion2() (version2.Database, error) {
	db := version2.NewDatabase(v.Path)

	err := db.SaveIndex(v.Index)
	if err != nil {
		return version2.Database{}, errors.Chain(err, "error saving index file")
	}

	return db, nil
}

func (v *Vault) ReloadIndex() error {
	var err error

	v.Index, err = v.Database.LoadIndex()
	if err != nil {
		return errors.Chain(err, "error loading index from database")
	}

	return nil
}
