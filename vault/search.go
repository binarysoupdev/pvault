package vault

import "slices"

func (v Vault) SearchNames(term string) []string {
	var matches []string

	if term == "" {
		matches = v.Map.GetNames()
	} else {
		matches = v.Map.SearchNames(term)
	}

	slices.Sort(matches)
	return matches
}
