package vault

import (
	"path/filepath"
	"pvault/vault/data"
	"pvault/vault/index"

	"github.com/google/uuid"
)

const CURRENT_VERSION uint16 = 2

type Vault struct {
	Path     string
	Index    index.IndexMap
	Database data.Database
}

func (v Vault) Version() uint16 {
	return v.Database.Version()
}

func (v Vault) RecordPath(id uuid.UUID) string {
	return filepath.Join(v.Path, id.String())
}
