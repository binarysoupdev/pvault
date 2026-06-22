package vault

import "slices"

func (v Vault) Search(term string) []string {
	var matches []string

	if term == "" {
		matches = v.Index.getNames()
	} else {
		matches = v.Index.searchNames(term)
	}

	slices.Sort(matches)
	return matches
}
