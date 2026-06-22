package flow

import (
	"flag"
	"pvault/vault"
	"strings"

	"github.com/binarysoupdev/got-style/style"
)

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
	f.display(v.Search(*f.term))
}

func (f SearchFlow) display(matches []string) {
	result := style.New(style.YELLOW)
	highlight := append(result, style.UNDERLINE)

	for i, match := range matches {
		start := strings.Index(strings.ToLower(match), strings.ToLower(*f.term))
		end := start + len(*f.term)

		result.Printf("[%d] %s", i+1, match[:start])
		highlight.Print(match[start:end])
		result.Println(match[end:])
	}
}
