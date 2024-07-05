package util

import (
	"fmt"
	"io"
	"os"
)

const progressbytes = `/-\|`

func DefaultCallback(w io.Writer) ParserCallback {
	var cnt int
	o, _ := os.Stdout.Stat()
	outtotty := (o.Mode()&os.ModeCharDevice) == os.ModeCharDevice || (o.Mode()&os.ModeNamedPipe) == os.ModeNamedPipe
	return func(s string, evt EVENT) error {
		w.Write(CreateOneLiner(s))
		if evt == ENDDOC {
			cnt++
			if !outtotty {
				if w == os.Stdout {
					fmt.Fprintf(os.Stderr, "\rdoc count: %d", cnt)
				} else {
					fmt.Fprintf(os.Stderr, "\r%c ", progressbytes[cnt%len(progressbytes)])
				}
			}
			w.Write([]byte{'\n'})
		}
		if evt == EOF && !outtotty && w != os.Stdout {
			fmt.Fprintf(os.Stderr, "\n\rdoc count: %d", cnt)
		}
		return nil
	}
}
