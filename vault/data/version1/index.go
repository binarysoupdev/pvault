package version1

import "path/filepath"

const INDEX_FILE = "index.txt"

func (db Database) IndexPath() string {
	return filepath.Join(db.Path, INDEX_FILE)
}
