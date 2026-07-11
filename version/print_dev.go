//go:build !prod

package version

import "github.com/binarysoupdev/got-style/style"

func Display() {
	style.New(style.BOLD, style.UNDERLINE).Printf("pvault (dev) @v%s\n", VERSION)
}
