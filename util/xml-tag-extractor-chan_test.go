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
	"testing"
)

//go:embed data/consolidated.xml
var filedata []byte

func TestParseXMLChanWithData(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	reader := bufio.NewReader(bytes.NewBuffer(filedata))
	datach := make(chan []byte)
	query := `CONSOLIDATED_LIST>INDIVIDUALS>INDIVIDUAL`
	errch := make(chan error)
	endl := []byte{}
	cb := DefaultCallback(writer)
	go func() { ParseXMLChan(datach, errch, query, cb) }()
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

	var err error

	for {
		e, ok := <-errch
		if !ok {
			break
		}
		err = fmt.Errorf("%w %w", err, e)
	}
	if err != nil && err != io.EOF {
		t.Errorf("failed with %v", err)
	}
	if writer.Len() < 1 {
		t.Errorf("no data written, check your source data and query")
	}
	fmt.Printf("%s\n", writer.String())
}

func TestParseXMLChan(t *testing.T) {
	xml := "<a>\n<b>hello\n</b>\n</a>"
	datach := make(chan []byte)
	errch := make(chan error)
	cb := DefaultCallback(os.Stdout)

	go func() { ParseXMLChan(datach, errch, "a>b", cb) }()
	datach <- []byte(xml)
	datach <- []byte{}
	close(datach)

	var err error

	for {
		e, ok := <-errch
		if !ok {
			break
		}
		err = fmt.Errorf("%w %w", err, e)
	}
	fmt.Printf("ERR: %v\n", err)
	if err != nil {
		t.Errorf("%v", err)
	}
}
