package flow

import (
	"flag"
	"pvault/vault"
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
		index: flags.Int("x", -1, "the index if multiple matches"),
	}
}

func (f SearchFlow) Display(v vault.Vault) {
	matches := v.Search(*f.term)

	if len(matches) == 0 {
		style.Info.Println("No MATCHES found")
		return
	}

	f.displayMatches(matches)
}

func (f SearchFlow) Select(v vault.Vault) string {
	matches := v.Search(*f.term)

	if len(matches) == 0 {
		style.Error.Println("No MATCHES found")
		return ""
	}

	if len(matches) == 1 {
		f.displayMatch(matches[0], 1)
		return matches[0]
	}

	index := *f.index - 1

	if index < 0 || index >= len(matches) {
		f.displayMatches(matches)
		style.Error.Println("Rerun with \"-x <index>\"")
		return ""
	}

	f.displayMatch(matches[index], index+1)
	return matches[index]
}

func (f SearchFlow) displayMatches(matches []string) {
	for i, match := range matches {
		f.displayMatch(match, i+1)
	}
}

func (f SearchFlow) displayMatch(match string, index int) {
	start := strings.Index(strings.ToLower(match), strings.ToLower(*f.term))
	end := start + len(*f.term)

	RESULT_STYLE.Printf("[%d] %s", index, match[:start])
	HIGHLIGHT_STYLE.Print(match[start:end])
	RESULT_STYLE.Println(match[end:])
}
