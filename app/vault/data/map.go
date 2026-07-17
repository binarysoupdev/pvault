package data

import (
	"strings"

	"github.com/google/uuid"
)

type NameMap map[string]uuid.UUID

func (m NameMap) FindName(id uuid.UUID) (string, bool) {
	for name, val := range m {
		if val == id {
			return name, true
		}
	}
	return "", false
}

func (m NameMap) GetNames() []string {
	names := make([]string, len(m))
	i := 0

	for name := range m {
		names[i] = name
		i++
	}

	return names
}

func (m NameMap) SearchNames(subStr string) []string {
	matches := []string{}

	for name := range m {
		if strings.Contains(strings.ToLower(name), strings.ToLower(subStr)) {
			matches = append(matches, name)
		}
	}

	return matches
}
