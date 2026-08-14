package v1

import (
	"bufio"
	"fmt"
	"io"
	"pvault/vault/index"
	"strings"

	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/google/uuid"
)

type Encoder struct{}

func (e Encoder) EncodeIndex(w io.Writer, idx index.IndexMap) error {
	entryNum := 0

	for name, id := range idx {
		_, err := fmt.Fprintf(w, "%s:%s\n", id.String(), name)
		if err != nil {
			return errors.Chain(err, fmt.Sprintf("error encoding entry [%d]", entryNum))
		}
		entryNum++
	}

	return nil
}

func (e Encoder) DecodeIndex(r io.Reader) (index.IndexMap, error) {
	idx := index.IndexMap{}
	line := 0

	scanner := bufio.NewScanner(r)
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
