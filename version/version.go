package version

import (
	"pvault/build"

	"github.com/binarysoupdev/got-style/style"
)

const VERSION = "2.0.0"

func Print() {
	style.New(style.BOLD, style.UNDERLINE).Printf("%s @v%s\n", build.AppName(), VERSION)
}
