package util

import (
	"fmt"
	"io"
	"os"
)

func DefaultCallback(w io.Writer) ParserCallback {
	var cnt int
	//var mutex sync.Mutex
	o, _ := os.Stdout.Stat()
	outtotty := (o.Mode()&os.ModeCharDevice) == os.ModeCharDevice || (o.Mode()&os.ModeNamedPipe) == os.ModeNamedPipe
	return func(s string, evt EVENT) error {
		//mutex.Lock()
		w.Write(CreateOneLiner(s))
		if evt == ENDDOC {
			cnt++
			if !outtotty {
				fmt.Fprintf(os.Stderr, "\rdoc count: %d", cnt)
			}
			w.Write([]byte{'\n'})
		}
		//mutex.Unlock()
		return nil
	}
}
