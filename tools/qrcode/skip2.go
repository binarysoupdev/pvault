package qrcode

import (
	"fmt"

	"github.com/binarysoupdev/go-commando/errors"

	qrcode "github.com/skip2/go-qrcode"
)

const BLOCK = '\u2588'

type Skip2Renderer struct{}

func (r Skip2Renderer) RenderToStdout(text string) error {
	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		return errors.Chain(err, "error creating qr code")
	}

	for _, row := range qr.Bitmap() {
		for _, bit := range row {
			if bit {
				fmt.Printf("%c%c", BLOCK, BLOCK)
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Println()
	}

	return nil
}
