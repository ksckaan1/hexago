package project

import (
	"io"
	"regexp"
)

type decolorizer struct {
	w io.Writer
}

func (d *decolorizer) Write(p []byte) (n int, err error) {
	decolorized := rgxANSI.ReplaceAll(p, nil)
	_, err = d.w.Write(decolorized)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var rgxANSI = regexp.MustCompile(ansi)
