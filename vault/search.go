package vault

import "slices"

func (v Vault) Search(term string) []string {
	var matches []string

	if term == "" {
		matches = v.Index.GetNames()
	} else {
		matches = v.Index.SearchNames(term)
	}

	slices.Sort(matches)
	return matches
}
