package vault

import (
	"path/filepath"

	"github.com/google/uuid"
)

func (v Vault) IndexPath() string {
	return filepath.Join(v.Path, INDEX_FILE)
}

func (v Vault) RecordPath(id uuid.UUID) string {
	return filepath.Join(v.Path, id.String()+".json")
}
