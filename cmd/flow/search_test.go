package flow_test

import (
	"flag"
	"pvault/cmd/flow"
	"pvault/vault"
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
		Index: vault.IndexMap{
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

func (s *SearchFlowSuite) TestDisplayNoResults() {
	//-- arrange
	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", "no match")
	s.Flow.Display(s.Vault)

	//-- assert
	s.Assert().Contains(out.ReadLine(), "No MATCHES found")
}

func (s *SearchFlowSuite) TestDisplayPrintResults() {
	//-- arrange
	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", "foo")
	s.Flow.Display(s.Vault)

	//-- assert
	line := out.ReadLine()
	s.Assert().Contains(line, "[1]")
	s.Assert().Contains(line, "Foo")

	line = out.ReadLine()
	s.Assert().Contains(line, "[2]")
	s.Assert().Contains(line, "Foo")
}

func (s *SearchFlowSuite) TestSelectNoResults() {
	//-- arrange
	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", "no match")
	res := s.Flow.Select(s.Vault)

	//-- assert
	s.Assert().Empty(res)
	s.Assert().Contains(out.ReadLine(), "No MATCHES found")
}

func (s *SearchFlowSuite) TestSelectOneResultResult() {
	//-- arrange
	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", "bar")
	res := s.Flow.Select(s.Vault)

	//-- assert
	s.Assert().Equal("bar1", res)
	s.Assert().Contains(out.ReadLine(), "Bar")
}

func (s *SearchFlowSuite) TestSelectManyResultsNoIndex() {
	//-- arrange
	out := pipe.OpenStdout(3)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", "foo")
	res := s.Flow.Select(s.Vault)

	//-- assert
	s.Assert().Empty(res)

	s.Assert().Contains(out.ReadLine(), "Foo")
	s.Assert().Contains(out.ReadLine(), "Foo")
	s.Assert().Contains(out.ReadLine(), "Rerun with \"-x <index>\"")
}

func (s *SearchFlowSuite) TestSelectManyResultsInvalidIndex() {
	//-- arrange
	out := pipe.OpenStdout(3)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", "foo", "-x", "3")
	res := s.Flow.Select(s.Vault)

	//-- assert
	s.Assert().Empty(res)

	s.Assert().Contains(out.ReadLine(), "Foo")
	s.Assert().Contains(out.ReadLine(), "Foo")
	s.Assert().Contains(out.ReadLine(), "Rerun with \"-x <index>\"")
}

func (s *SearchFlowSuite) TestSelectManyResultsValidIndexReturnsMatch() {
	//-- arrange
	out := pipe.OpenStdout(2)
	defer out.Close()

	//-- act
	s.ParseFlags("-s", "foo", "-x", "1")
	res := s.Flow.Select(s.Vault)

	//-- assert
	s.Assert().Equal("foo1", res)
	s.Assert().Contains(out.ReadLine(), "Foo")
}
