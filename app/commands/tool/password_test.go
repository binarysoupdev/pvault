package tool_test

import (
	"fmt"
	cmd "pvault/app/commands/tool"
	"pvault/tools/clipboard"
	"testing"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/go-commando/test"
	"github.com/binarysoupdev/tinsel/pipe"
	"github.com/stretchr/testify/suite"
)

type PasswordTestSuite struct {
	test.CommandSuite[*cmd.PasswordCommand]
	Clipboard *clipboard.MockClipboard
}

func TestCopyCommandSuite(t *testing.T) {
	s := PasswordTestSuite{
		Clipboard: clipboard.Mock(),
	}

	s.CommandSuite = test.NewCommandSuite(cmd.NewPasswordCommand(s.Clipboard))
	suite.Run(t, &s)
}

func (s *PasswordTestSuite) SetupTest() {
	*s.Clipboard = clipboard.MockClipboard{}
}

//=====================================

func (s *PasswordTestSuite) TestRunFailsWhenLengthTooShort() {
	//-- act
	s.RunCommand("-len", "0")

	//-- assert
	s.RequireResultFail("\"len\" too short")
}

func (s *PasswordTestSuite) TestRunPassesAndPrintsPassword() {
	//-- arrange
	const LENGTH = 50

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-len", fmt.Sprintf("%d", LENGTH))

	//-- assert
	s.RequireResultPass()
	s.Assert().Len(out.ReadLine(), LENGTH)
}

func (s *PasswordTestSuite) TestRunCopyFailsWhenClipboardUnsupported() {
	//-- arrange
	s.Clipboard.UnsupportedError = errors.New("")

	//-- act
	s.RunCommand("-copy")

	//-- assert
	s.RequireResultFail("clipboard unsupported")
}

func (s *PasswordTestSuite) TestRunCopyFailsWhenErrorCopyingToClipboard() {
	//-- arrange
	s.Clipboard.WriteError = errors.New("")

	//-- act
	s.RunCommand("-copy")

	//-- assert
	s.RequireResultFail("error copying to clipboard")
}

func (s *PasswordTestSuite) TestRunCopyPassesAndCopiesPasswordToClipboard() {
	//-- arrange
	const LENGTH = 50

	out := pipe.OpenStdout(1)
	defer out.Close()

	//-- act
	s.RunCommand("-copy", "-len", fmt.Sprintf("%d", LENGTH))

	//-- assert
	s.RequireResultPass()
	s.Assert().Contains(out.ReadLine(), "[+] PASSWORD copied to clipboard")

	s.Assert().Len(s.Clipboard.Data, LENGTH)
}
