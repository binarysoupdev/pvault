package v1

import (
	"fmt"
	"os"
	"pvault/vault/index"

	"github.com/binarysoupdev/go-commando/errors"
)

func (idx Index) SaveIndex(m index.IndexMap) error {
	file, err := os.Create(idx.filepath())
	if err != nil {
		return errors.Chain(err, "error creating index file")
	}
	defer file.Close()

	for name, id := range m {
		fmt.Fprintf(file, "%s:%s\n", id.String(), name)
	}

	return nil
}
