package search_test

import (
	"flag"
	"fmt"
	vault_flow "pvault/app/flow/vault"
	search_flow "pvault/app/flow/vault/search"
	"testing"

	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/suite"
)

type SearchFlowSuite struct {
	suite.Suite
	VaultMock *vault_flow.VaultMock

	Flow  search_flow.SearchFlow
	flags *flag.FlagSet
}

func TestSearchFlowSuite(t *testing.T) {
	suite.Run(t, &SearchFlowSuite{})
}

func (s *SearchFlowSuite) SetupTest() {
	s.VaultMock = &vault_flow.VaultMock{}

	s.flags = flag.NewFlagSet("", flag.PanicOnError)
	s.Flow = search_flow.NewSearchFlow(s.flags)
}

func (s *SearchFlowSuite) ParseFlags(args ...string) {
	s.flags.Parse(args)
}

//=====================================

func (s *SearchFlowSuite) TestDisplayReturnsErrorWhenNoResults() {
	//-- arrange
	const SEARCH = "term"
	s.VaultMock.SearchResults = []string{}

	//-- act
	s.ParseFlags("-s", SEARCH)
	res := s.Flow.Display(s.VaultMock)

	//-- assert
	s.Assert().ErrorContains(res, "no matches found")
	s.Assert().Equal(SEARCH, s.VaultMock.SearchTermParam)
}

func (s *SearchFlowSuite) TestDisplayReturnsNoErrorAndPrintResults() {
	//-- arrange
	const SEARCH = "Foo"
	s.VaultMock.SearchResults = []string{"Foo1", "Foo2"}

	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", SEARCH)
	err := s.Flow.Display(s.VaultMock)

	//-- assert
	s.Require().NoError(err)
	s.Assert().Equal(SEARCH, s.VaultMock.SearchTermParam)

	line := out.ReadLine()
	s.Assert().Contains(line, "[0]")
	s.Assert().Contains(line, SEARCH)

	line = out.ReadLine()
	s.Assert().Contains(line, "[1]")
	s.Assert().Contains(line, SEARCH)
}

func (s *SearchFlowSuite) TestSelectReturnsErrorWhenNoResults() {
	//-- arrange
	const SEARCH = "term"
	s.VaultMock.SearchResults = []string{}

	//-- act
	s.ParseFlags("-s", "term")
	_, res := s.Flow.Select(s.VaultMock)

	//-- assert
	s.Assert().ErrorContains(res, "no matches found")
	s.Assert().Equal(SEARCH, s.VaultMock.SearchTermParam)
}

func (s *SearchFlowSuite) TestSelectReturnsNoErrorAndMatchWhenOnlyOneResult() {
	//-- arrange
	const SEARCH = "Foo"
	s.VaultMock.SearchResults = []string{"Foo1"}

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", SEARCH)
	res, err := s.Flow.Select(s.VaultMock)

	//-- assert
	s.Require().NoError(err)
	s.Assert().Equal(SEARCH, s.VaultMock.SearchTermParam)
	s.Assert().Equal(s.VaultMock.SearchResults[0], res)

	line := out.ReadLine()
	s.Assert().Contains(line, "[0]")
	s.Assert().Contains(line, SEARCH)
}

func (s *SearchFlowSuite) TestSelectReturnsErrorWhenManyResultsAndInvalidIndex() {
	//-- arrange
	const SEARCH = "Foo"
	const INDEX = 2
	s.VaultMock.SearchResults = []string{"Foo1", "Foo2"}

	io := pipe.OpenStdio(1, 3, false)
	defer io.Close()

	io.Queue("INDEX: ", INDEX)
	io.EndQueue()

	//-- act
	s.ParseFlags("-s", SEARCH)
	_, res := s.Flow.Select(s.VaultMock)

	//-- assert
	s.Assert().ErrorContains(res, fmt.Sprintf("invalid index \"%d\"", INDEX))
	s.Assert().Equal(SEARCH, s.VaultMock.SearchTermParam)

	s.Assert().Contains(io.ReadLine(), "[0]")
	s.Assert().Contains(io.ReadLine(), "[1]")
	s.Assert().Contains(io.ReadLine(), "Enter INDEX")
}

func (s *SearchFlowSuite) TestSelectReturnsMatchAndNoErrorWhenManyResultsAndValidIndex() {
	//-- arrange
	const SEARCH = "Foo"
	const INDEX = 1
	s.VaultMock.SearchResults = []string{"Foo1", "Foo2"}

	io := pipe.OpenStdio(1, 3, false)
	defer io.Close()

	io.Queue("INDEX: ", INDEX)
	io.EndQueue()

	//-- act
	s.ParseFlags("-s", SEARCH)
	res, err := s.Flow.Select(s.VaultMock)

	//-- assert
	s.Require().NoError(err)
	s.Assert().Equal(SEARCH, s.VaultMock.SearchTermParam)
	s.Assert().Equal(s.VaultMock.SearchResults[INDEX], res)

	s.Assert().Contains(io.ReadLine(), "[0]")
	s.Assert().Contains(io.ReadLine(), "[1]")
	s.Assert().Contains(io.ReadLine(), "Enter INDEX")
}
