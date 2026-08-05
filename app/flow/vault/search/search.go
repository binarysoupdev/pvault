package search

import (
	"errors"
	"flag"
	"pvault/app/flow/prompt"
	vault_flow "pvault/app/flow/vault"
	"strings"

	"github.com/binarysoupdev/got-style/style"
)

var RESULT_STYLE = style.New(style.YELLOW)
var HIGHLIGHT_STYLE = style.New(style.YELLOW, style.UNDERLINE)

type SearchFlow struct {
	term  *string
	index *int
}

func NewSearchFlow(flags *flag.FlagSet) SearchFlow {
	return SearchFlow{
		term:  flags.String("s", "", "the search term"),
		index: flags.Int("x", -1, "the index if many matches"),
	}
}

func (f SearchFlow) Display(v vault_flow.Vault) error {
	matches := v.SearchNames(*f.term)

	if len(matches) == 0 {
		return errors.New("no matches found")
	}

	f.displayMatches(matches)
	return nil
}

func (f SearchFlow) Select(v vault_flow.Vault) (string, error) {
	matches := v.SearchNames(*f.term)

	if len(matches) == 0 {
		return "", errors.New("no matches found")
	}

	if len(matches) == 1 {
		f.displayMatch(matches[0], 0)
		return matches[0], nil
	}

	f.displayMatches(matches)
	index := f.selectIndex(len(matches))

	return matches[index], nil
}

func (f SearchFlow) displayMatches(matches []string) {
	for i, match := range matches {
		f.displayMatch(match, i)
	}
}

func (f SearchFlow) displayMatch(match string, index int) {
	start := strings.Index(strings.ToLower(match), strings.ToLower(*f.term))
	end := start + len(*f.term)

	RESULT_STYLE.Printf("[%d] %s", index, match[:start])
	HIGHLIGHT_STYLE.Print(match[start:end])
	RESULT_STYLE.Println(match[end:])
}

func (f SearchFlow) selectIndex(numMatches int) int {
	index := *f.index

	for index < 0 || index >= numMatches {
		index = prompt.Number("Select MATCH: ", -1)
	}

	return index
}
