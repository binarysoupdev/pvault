package vault

import (
	"encoding/json"
	"os"
	"path/filepath"
	"pvault/chain"

	"github.com/google/uuid"
)

const VAULT_DIR = "tmp/vault"

type Vault struct{}

func (v Vault) NewRecord(name string) (Record, error) {
	r := Record{
		ID:       uuid.New(),
		Name:     name,
		Username: "",
		Password: "",
		Other:    map[string]string{},
	}

	return r, v.SaveRecord(r)
}

func (Vault) SaveRecord(r Record) error {
	//TODO: check record name is unique

	file, err := os.Create(filepath.Join(VAULT_DIR, r.ID.String()+".json"))
	if err != nil {
		return chain.Error(err, "error creating record file")
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(r)
	if err != nil {
		return chain.Error(err, "error encoding record JSON")
	}

	return nil
}
