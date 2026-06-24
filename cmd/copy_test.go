package cmd_test

import (
	"errors"
	"pvault/cmd"
	"pvault/tools/clipboard"
	"testing"

	"github.com/binarysoupdev/go-commando/test"
	"github.com/stretchr/testify/suite"
)

type CopyTestSuite struct {
	test.CommandSuite[*cmd.CopyCommand]
	Clipboard *clipboard.MockClipboard
}

func TestCopyCommandSuite(t *testing.T) {
	s := CopyTestSuite{
		Clipboard: clipboard.Mock(),
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewCopyCommand(s.Clipboard))
	suite.Run(t, &s)
}

func (s *CopyTestSuite) SetupTest() {
	*s.Clipboard = clipboard.MockClipboard{}
}

//=====================================

func (s *CopyTestSuite) TestRunFailClipboardUnsupported() {
	//-- arrange
	s.Clipboard.Unsupported = errors.New("")

	//-- act
	s.RunCommand()

	//-- assert
	s.RequireResultFail("clipboard unsupported")
}
