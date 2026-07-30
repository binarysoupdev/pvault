package vault

import (
	"path/filepath"
)

const METADATA_FILE = "VAULT"

func metadataPath(path string) string {
	return filepath.Join(path, METADATA_FILE)
}
