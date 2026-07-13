package version1

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"pvault/vault/index"
	"strings"

	"github.com/binarysoupdev/go-commando/errors"

	"github.com/google/uuid"
)

const INDEX_FILE = "index.txt"

func (db Database) IndexPath() string {
	return filepath.Join(db.Path, INDEX_FILE)
}

func (db Database) SaveIndex(idx index.IndexMap) error {
	file, err := os.Create(db.IndexPath())
	if err != nil {
		return errors.Chain(err, "error creating index file")
	}
	defer file.Close()

	for name, id := range idx {
		fmt.Fprintf(file, "%s:%s\n", id.String(), name)
	}

	return nil
}

func (db Database) LoadIndex() (index.IndexMap, error) {
	file, err := os.Open(db.IndexPath())
	if err != nil {
		return nil, errors.Chain(err, "error opening index file")
	}
	defer file.Close()

	idx := index.IndexMap{}
	line := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line++

		tokens := strings.SplitN(scanner.Text(), ":", 2)
		if len(tokens) < 2 {
			return idx, errors.Format("[line %d] invalid index pair", line)
		}

		name := tokens[1]
		id, err := uuid.Parse(tokens[0])
		if err != nil {
			return idx, errors.ChainFormat(err, "[line %d] invalid uuid", line)
		}

		idx[name] = id
	}

	return idx, scanner.Err()
}
