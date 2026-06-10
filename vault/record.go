package vault

import (
	"encoding/json"
	"os"
	"pvault/chain"

	"github.com/google/uuid"
)

type Record struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	Username string `json:"username"`
	Password string `json:"password"`
	Other    any    `json:"other"`
}

func SaveRecordJSON(r Record, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return chain.Error(err, "error creating record file")
	}
	defer file.Close()

	e := json.NewEncoder(file)
	e.SetIndent("", "    ")

	err = e.Encode(r)
	if err != nil {
		return chain.Error(err, "error encoding record JSON")
	}

	return nil
}

func LoadRecordJSON(path string) (Record, error) {
	var r Record

	file, err := os.Open(path)
	if err != nil {
		return r, chain.Error(err, "error opening record file")
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&r)
	if err != nil {
		return r, chain.Error(err, "error decoding record JSON")
	}

	return r, nil
}
