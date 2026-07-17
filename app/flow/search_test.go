package flow_test

import (
	"flag"
	"pvault/app/flow"
	"pvault/app/vault"
	"pvault/app/vault/data"
	"testing"

	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type SearchFlowSuite struct {
	suite.Suite
	Vault vault.Vault

	Flow  flow.SearchFlow
	flags *flag.FlagSet
}

func TestSearchFlowSuite(t *testing.T) {
	suite.Run(t, &SearchFlowSuite{})
}

func (s *SearchFlowSuite) SetupTest() {
	s.Vault = vault.Vault{
		Index: data.NameMap{
			"Foo1": uuid.New(),
			"Foo2": uuid.New(),
			"Bar1": uuid.New(),
		},
	}

	s.flags = flag.NewFlagSet("", flag.PanicOnError)
	s.Flow = flow.NewSearchFlow(s.flags)
}

func (s *SearchFlowSuite) ParseFlags(args ...string) {
	s.flags.Parse(args)
}

//=====================================

func (s *SearchFlowSuite) TestDisplayReturnsErrorWhenNoResults() {
	//-- act
	s.ParseFlags("-s", "no match")
	res := s.Flow.Display(s.Vault)

	//-- assert
	s.Assert().ErrorContains(res, "no matches found")
}

func (s *SearchFlowSuite) TestDisplayReturnsNoErrorAndPrintResults() {
	//-- arrange
	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", "foo")
	err := s.Flow.Display(s.Vault)

	//-- assert
	s.Require().NoError(err)

	line := out.ReadLine()
	s.Assert().Contains(line, "[1]")
	s.Assert().Contains(line, "Foo")

	line = out.ReadLine()
	s.Assert().Contains(line, "[2]")
	s.Assert().Contains(line, "Foo")
}

func (s *SearchFlowSuite) TestSelectReturnsErrorWhenNoResults() {
	//-- act
	s.ParseFlags("-s", "no match")
	_, res := s.Flow.Select(s.Vault)

	//-- assert
	s.Assert().ErrorContains(res, "no matches found")
}

func (s *SearchFlowSuite) TestSelectReturnsNoErrorAndMatchWhenOnlyOneResult() {
	//-- arrange
	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", "bar")
	res, err := s.Flow.Select(s.Vault)

	//-- assert
	s.Require().NoError(err)

	s.Assert().Equal("Bar1", res)
	s.Assert().Contains(out.ReadLine(), "Bar")
}

func (s *SearchFlowSuite) TestSelectReturnsErrorWhenManyResultsAndNoIndex() {
	//-- arrange
	out := pipe.OpenStdout(3)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", "foo")
	_, res := s.Flow.Select(s.Vault)

	//-- assert
	s.Assert().ErrorContains(res, "rerun with \"-x <index>\"")

	s.Assert().Contains(out.ReadLine(), "Foo")
	s.Assert().Contains(out.ReadLine(), "Foo")
}

func (s *SearchFlowSuite) TestSelectReturnsErrorWhenManyResultsAndInvalidIndex() {
	//-- arrange
	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", "foo", "-x", "3")
	_, res := s.Flow.Select(s.Vault)

	//-- assert
	s.Assert().ErrorContains(res, "rerun with \"-x <index>\"")

	s.Assert().Contains(out.ReadLine(), "Foo")
	s.Assert().Contains(out.ReadLine(), "Foo")
}

func (s *SearchFlowSuite) TestSelectReturnsNoErrorAndMatchWhenManyResultsAndValidIndex() {
	//-- arrange
	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", "foo", "-x", "1")
	res, err := s.Flow.Select(s.Vault)

	//-- assert
	s.Require().NoError(err)

	s.Assert().Equal("Foo1", res)
	s.Assert().Contains(out.ReadLine(), "Foo")
}
