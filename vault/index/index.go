package index

import (
	"github.com/google/uuid"
)

const (
	VERSION           = 1
	INDEX_FILE        = "index.bin"
	LEGACY_INDEX_FILE = "index.txt"
)

type IndexMap map[string]uuid.UUID
