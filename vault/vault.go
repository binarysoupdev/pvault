package vault

import (
	"os"
	"path/filepath"
	"pvault/chain"
)

const INDEX_FILE = "index.bin"

type Vault struct {
	Index IndexMap
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

	err = IndexMap{}.Save(filepath.Join(path, INDEX_FILE))
	if err != nil {
		return chain.Error(err, "error saving index file")
	}

	return nil
}
