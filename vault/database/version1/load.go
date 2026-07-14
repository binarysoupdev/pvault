package v1

import (
	"bufio"
	"os"
	"pvault/vault/index"
	"strings"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func (idx Index) LoadIndex() (index.IndexMap, error) {
	file, err := os.Open(idx.Filepath())
	if err != nil {
		return nil, errors.Chain(err, "error opening index file")
	}
	defer file.Close()

	m := index.IndexMap{}
	line := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line++

		tokens := strings.SplitN(scanner.Text(), ":", 2)
		if len(tokens) < 2 {
			return m, errors.Format("[line %d] invalid index pair", line)
		}

		name := tokens[1]
		id, err := uuid.Parse(tokens[0])
		if err != nil {
			return m, errors.ChainFormat(err, "[line %d] invalid uuid", line)
		}

		m[name] = id
	}

	return m, scanner.Err()
}
