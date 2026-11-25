package clipboard

import (
	"bytes"
	"os/exec"

	"github.com/binarysoupdev/go-extensions/errors"
)

const xclip = "xclip"

type XClipClipboard struct{}

func (XClipClipboard) CheckUnsupported() error {
	if _, err := exec.LookPath(xclip); err != nil {
		return errors.New("xclip not found")
	}
	return nil
}

func (XClipClipboard) Read() (string, error) {
	cmd := exec.Command(xclip, "-out", "-selection", "clipboard")

	out, err := cmd.Output()
	if err != nil {
		return "", errors.Chain(err, "error running xclip")
	}

	return string(out), nil
}

func (XClipClipboard) Write(data string) error {
	buf := bytes.NewBuffer([]byte(data))

	cmd := exec.Command(xclip, "-in", "-selection", "clipboard")
	cmd.Stdin = buf

	if err := cmd.Run(); err != nil {
		return errors.Chain(err, "error running xclip")
	}
	return nil
}
