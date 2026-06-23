package vault

import (
	"os"
	"path/filepath"
	"pvault/chain"
	"pvault/vault/index"

	"github.com/google/uuid"
)

const INDEX_FILE = "index.bin"

type Vault struct {
	Path  string
	Index index.IndexMap
}

func InitializeNew(path string) error {
	_, err := os.Stat(path)
	if err == nil || !os.IsNotExist(err) {
		return chain.New("vault path already exists")
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return chain.Error(err, "error creating vault directory")
	}

	err = index.IndexMap{}.Save(filepath.Join(path, INDEX_FILE))
	if err != nil {
		return chain.Error(err, "error saving index file")
	}

	return nil
}

func Open(path string) (Vault, error) {
	v := Vault{
		Path: path,
	}

	var err error
	v.Index, err = index.LoadIndex(v.IndexPath())
	if err != nil {
		return v, chain.Error(err, "error loading index file")
	}

	return v, nil
}

func (v Vault) IndexPath() string {
	return filepath.Join(v.Path, INDEX_FILE)
}

func (v Vault) RecordPath(id uuid.UUID) string {
	return filepath.Join(v.Path, id.String()+".json")
}
