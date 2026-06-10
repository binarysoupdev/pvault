package data

import (
	"encoding/json"
	"os"
)

func LoadJSON[T any](path string) (T, error) {
	var obj T

	file, err := os.Open(path)
	if err != nil {
		return obj, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&obj)
	if err != nil {
		return obj, err
	}

	return obj, nil
}

func SaveJSON[T any](obj T, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	e := json.NewEncoder(file)
	e.SetIndent("", "    ")

	err = e.Encode(obj)
	if err != nil {
		return err
	}

	return nil
}
