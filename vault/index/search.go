package index

import (
	"strings"

	"github.com/google/uuid"
)

func (idx IndexMap) FindName(id uuid.UUID) (string, bool) {
	for name, val := range idx {
		if val == id {
			return name, true
		}
	}
	return "", false
}

func (idx IndexMap) GetNames() []string {
	names := make([]string, len(idx))
	i := 0

	for name := range idx {
		names[i] = name
		i++
	}

	return names
}

func (idx IndexMap) SearchNames(subStr string) []string {
	matches := []string{}

	for name := range idx {
		if strings.Contains(strings.ToLower(name), strings.ToLower(subStr)) {
			matches = append(matches, name)
		}
	}

	return matches
}
