package vault

import (
	"path/filepath"
	"pvault/vault/data"
	"pvault/vault/index"

	"github.com/google/uuid"
)

const (
	CURRENT_VERSION uint16 = 2
	INDEX_FILE             = "index.bin"
)

type Vault struct {
	Path    string
	Version uint16

	Database data.Database
	Index    index.IndexMap
}

func (v Vault) IndexPath() string {
	return filepath.Join(v.Path, INDEX_FILE)
}

func (v Vault) RecordPath(id uuid.UUID) string {
	return filepath.Join(v.Path, id.String())
}
