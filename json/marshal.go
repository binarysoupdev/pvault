package json

import (
	"encoding/json"
	"os"
	"pvault/errors"
)

func Marshal[T any](obj T) ([]byte, error) {
	return json.Marshal(obj)
}

func MarshalFile[T any](obj T, path string) error {
	return MarshalFilePretty(obj, path, "")
}

func MarshalFilePretty[T any](obj T, path string, indent string) error {
	bytes, err := json.MarshalIndent(obj, "", indent)
	if err != nil {
		return errors.Chain(err, "error marshaling JSON")
	}

	err = os.WriteFile(path, bytes, 0666)
	if err != nil {
		return errors.Chain(err, "error writing JSON file")
	}

	return nil
}
