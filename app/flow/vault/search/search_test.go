package search_test

import (
	"flag"
	"pvault/app/flow"
	"pvault/app/vault"
	"testing"

	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/suite"
)

type SearchFlowSuite struct {
	suite.Suite
	Vault *vault.Mock

	Flow  flow.SearchFlow
	flags *flag.FlagSet
}

func TestSearchFlowSuite(t *testing.T) {
	suite.Run(t, &SearchFlowSuite{})
}

func (s *SearchFlowSuite) SetupTest() {
	s.Vault = &vault.Mock{}

	s.flags = flag.NewFlagSet("", flag.PanicOnError)
	s.Flow = flow.NewSearchFlow(s.flags)
}

func (s *SearchFlowSuite) ParseFlags(args ...string) {
	s.flags.Parse(args)
}

//=====================================

func (s *SearchFlowSuite) TestDisplayReturnsErrorWhenNoResults() {
	//-- arrange
	const SEARCH = "term"
	s.Vault.SearchResults = []string{}

	//-- act
	s.ParseFlags("-s", SEARCH)
	res := s.Flow.Display(s.Vault)

	//-- assert
	s.Assert().ErrorContains(res, "no matches found")
	s.Assert().Equal(SEARCH, s.Vault.SearchTermParam)
}

func (s *SearchFlowSuite) TestDisplayReturnsNoErrorAndPrintResults() {
	//-- arrange
	const SEARCH = "Foo"
	s.Vault.SearchResults = []string{"Foo1", "Foo2"}

	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", SEARCH)
	err := s.Flow.Display(s.Vault)

	//-- assert
	s.Require().NoError(err)
	s.Assert().Equal(SEARCH, s.Vault.SearchTermParam)

	line := out.ReadLine()
	s.Assert().Contains(line, "[1]")
	s.Assert().Contains(line, SEARCH)

	line = out.ReadLine()
	s.Assert().Contains(line, "[2]")
	s.Assert().Contains(line, SEARCH)
}

func (s *SearchFlowSuite) TestSelectReturnsErrorWhenNoResults() {
	//-- arrange
	const SEARCH = "term"
	s.Vault.SearchResults = []string{}

	//-- act
	s.ParseFlags("-s", "term")
	_, res := s.Flow.Select(s.Vault)

	//-- assert
	s.Assert().ErrorContains(res, "no matches found")
	s.Assert().Equal(SEARCH, s.Vault.SearchTermParam)
}

func (s *SearchFlowSuite) TestSelectReturnsNoErrorAndMatchWhenOnlyOneResult() {
	//-- arrange
	const SEARCH = "Foo"
	s.Vault.SearchResults = []string{"Foo1"}

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", SEARCH)
	res, err := s.Flow.Select(s.Vault)

	//-- assert
	s.Require().NoError(err)
	s.Assert().Equal(SEARCH, s.Vault.SearchTermParam)

	s.Assert().Equal(s.Vault.SearchResults[0], res)
	s.Assert().Contains(out.ReadLine(), SEARCH)
}

func (s *SearchFlowSuite) TestSelectReturnsErrorWhenManyResultsAndNoIndex() {
	//-- arrange
	const SEARCH = "Foo"
	s.Vault.SearchResults = []string{"Foo1", "Foo2"}

	out := pipe.OpenStdout(3)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", SEARCH)
	_, res := s.Flow.Select(s.Vault)

	//-- assert
	s.Assert().ErrorContains(res, "rerun with \"-x <index>\"")
	s.Assert().Equal(SEARCH, s.Vault.SearchTermParam)

	s.Assert().Contains(out.ReadLine(), SEARCH)
	s.Assert().Contains(out.ReadLine(), SEARCH)
}

func (s *SearchFlowSuite) TestSelectReturnsErrorWhenManyResultsAndInvalidIndex() {
	//-- arrange
	const SEARCH = "Foo"
	s.Vault.SearchResults = []string{"Foo1", "Foo2"}

	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", SEARCH, "-x", "3")
	_, res := s.Flow.Select(s.Vault)

	//-- assert
	s.Assert().ErrorContains(res, "rerun with \"-x <index>\"")
	s.Assert().Equal(SEARCH, s.Vault.SearchTermParam)

	s.Assert().Contains(out.ReadLine(), SEARCH)
	s.Assert().Contains(out.ReadLine(), SEARCH)
}

func (s *SearchFlowSuite) TestSelectReturnsMatchAndNoErrorWhenManyResultsAndValidIndex() {
	//-- arrange
	const SEARCH = "Foo"
	s.Vault.SearchResults = []string{"Foo1", "Foo2"}

	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", SEARCH, "-x", "1")
	res, err := s.Flow.Select(s.Vault)

	//-- assert
	s.Require().NoError(err)
	s.Assert().Equal(SEARCH, s.Vault.SearchTermParam)

	s.Assert().Equal(s.Vault.SearchResults[0], res)
	s.Assert().Contains(out.ReadLine(), SEARCH)
}
