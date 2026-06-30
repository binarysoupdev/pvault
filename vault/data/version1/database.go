package version1

type Database struct {
	Path string
}

func NewDatabase(path string) Database {
	return Database{
		Path: path,
	}
}

func (Database) GetVersion() uint16 {
	return 1
}
