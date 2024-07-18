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
	oneliner := opts.MakeOneLiner
	return func(s string, evt EVENT) error {
		w.Write(createoneliner(s, oneliner))
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
		if evt == EOF {
			if !outtotty && w != os.Stdout {
				fmt.Fprintf(os.Stderr, "\n\n\rdoc count: %d", cnt)
			}
		}
		return nil
	}
}

func createoneliner(s string, doit bool) []byte {
	if doit {
		return CreateOneLiner(s)
	}
	return []byte(s)
}
