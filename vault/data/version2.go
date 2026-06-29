package data

type DatabaseV2 struct {
	Path string
}

func NewDatabaseV2(path string) DatabaseV2 {
	return DatabaseV2{
		Path: path,
	}
}

func (DatabaseV2) Version() uint16 {
	return 2
}

func (DatabaseV2) Upgrade() error {
	return nil
}
