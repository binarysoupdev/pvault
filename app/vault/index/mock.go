package index

import (
	"path/filepath"
	"pvault/app/vault/data"
	v2 "pvault/app/vault/index/version2"

	"github.com/google/uuid"
)

const MOCK_INDEX_FILE = "mock_index"

type Mock struct {
	Version      int
	Path         string
	Map          data.NameMap
	UpgradeIndex v2.Index

	SaveMapError error
	LoadMapError error
	UpgradeError error
}

func (m Mock) GetVersion() int {
	return m.Version
}

func (m Mock) Filepath() string {
	return filepath.Join(m.Path, MOCK_INDEX_FILE)
}

func (m Mock) RecordPath(id uuid.UUID) string {
	return filepath.Join(m.Path, id.String())
}

func (m *Mock) SaveMap(nameMap data.NameMap) error {
	m.Map = nameMap
	return m.SaveMapError
}

func (m Mock) LoadMap() (data.NameMap, error) {
	return m.Map, m.LoadMapError
}

func (m Mock) Upgrade() (v2.Index, error) {
	return m.UpgradeIndex, m.UpgradeError
}
