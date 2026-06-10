package cmd_test

import (
	"pvault/cmd"
	"testing"

	"github.com/binarysoupdev/go-commando/test"
	"github.com/stretchr/testify/suite"
)

type CreateTestSuite struct {
	test.CommandSuite[*cmd.CreateCommand]
}

func TestCreateCommandSuite(t *testing.T) {
	suite.Run(t, &CreateTestSuite{
		CommandSuite: test.NewCommandSuite(cmd.NewCreateCommand()),
	})
}

func (s *CreateTestSuite) TestNameNotEmpty() {
	//-- act
	s.RunCommand("-name", "")

	//-- assert
	s.RequireResultFail("\"name\" cannot be empty")
}

func (s *CreateTestSuite) TestCreateRecord() {
	//TODO: need way to set vault/local directory
}
