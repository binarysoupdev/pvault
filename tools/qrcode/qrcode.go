package qrcode

type Renderer interface {
	RenderToStdout(text string) error
}
