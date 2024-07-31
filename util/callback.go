package util

import (
	"fmt"
	"io"
	"os"
)

const progressbytes = `/-\|`

func DefaultCallback(w io.Writer, opts *Options) ParserCallback {
	var cnt int
	o, _ := os.Stdout.Stat()
	outtotty := (o.Mode()&os.ModeCharDevice) == os.ModeCharDevice || (o.Mode()&os.ModeNamedPipe) == os.ModeNamedPipe
	return func(s string, evt EVENT) error {
		if _, err := w.Write([]byte(s)); err != nil {
			return err
		}
		if evt == ENDDOC {
			cnt++
			if !outtotty {
				if w == os.Stdout {
					fmt.Fprintf(os.Stderr, "\rdoc count: %d", cnt)
				} else {
					fmt.Fprintf(os.Stderr, "\r%c ", progressbytes[cnt%len(progressbytes)])
				}
			}
			if _, err := w.Write([]byte{'\n'}); err != nil {
				return err
			}
		}
		if evt == EOF {
			if !outtotty && w != os.Stdout {
				fmt.Fprintf(os.Stderr, "\n\n\rdoc count: %d", cnt)
			}
		}
		return nil
	}
}
