package util

// go test -v -count=1 -run ^TestParseXMLChan$ ./util
// go test -v -count=1 -run ^TestParseXMLChanWithData$ ./util

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
)

//go:embed data/consolidated.xml
var filedata []byte

func TestParseXMLChanWithData(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	reader := bufio.NewReader(bytes.NewBuffer(filedata))
	query := `CONSOLIDATED_LIST>INDIVIDUALS>INDIVIDUAL`
	datach := make(chan []byte)
	errch := make(chan error)
	endl := []byte{}
	opts := NewOpts()
	cb := DefaultCallback(writer, opts)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() { defer wg.Done(); ParseXMLChan(datach, errch, query, cb) }()
	cnt := 0
	for {
		s, err := reader.ReadString('\n')
		cnt++
		if err != nil {
			fmt.Printf("line# %d loop err: %v\n", cnt, err)
			break
		}
		//println("sending", s)
		datach <- []byte(s)
		datach <- endl
	}
	close(datach)

	go func() {
		for e := range errch {
			fmt.Fprintf(os.Stderr, ">> %v\n", e)
			if e != nil && e != io.EOF {
				t.Errorf("failed with %v", e)
			}
		}
		wg.Wait()
		close(errch)
	}()

	if writer.Len() < 1 {
		t.Errorf("no data written, check your source data and query")
	}
	fmt.Printf("%s\n", writer.String())
}

func TestParseXMLChan(t *testing.T) {
	xml := "<a>\n<b>hello\n</b>\n</a>"
	datach := make(chan []byte)
	errch := make(chan error)
	cb := DefaultCallback(os.Stdout, NewOpts())
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() { defer wg.Done(); ParseXMLChan(datach, errch, "a>b", cb) }()
	datach <- []byte(xml)
	datach <- []byte{}
	close(datach)

	go func() {
		for e := range errch {
			if e != nil && e != io.EOF {
				t.Errorf("failed with %v", e)
			}
		}
		wg.Wait()
		close(errch)
	}()
}
