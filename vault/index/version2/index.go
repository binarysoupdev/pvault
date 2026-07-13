package v2

import (
	"path/filepath"
	"pvault/common"
	"pvault/vault/index"

	"github.com/google/uuid"
)

const (
	VERSION    = 2
	INDEX_FILE = "index.bin"
)

type Index struct {
	Path string
}

func NewIndex(path string) Index {
	return Index{
		Path: path,
	}
}

func (idx Index) filepath() string {
	return filepath.Join(idx.Path, INDEX_FILE)
}

func (Index) GetVersion() uint16 {
	return VERSION
}

func (idx Index) RecordPath(id uuid.UUID) string {
	return filepath.Join(idx.Path, id.String())
}

func (Index) Upgrade(m index.IndexMap) error {
	return common.NotSupportedError{}
}
