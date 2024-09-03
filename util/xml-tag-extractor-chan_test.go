package util

// go test -v -count=1 -run ^TestParseXMLChan$ ./util
// go test -v -count=1 -run ^TestParseXMLChanWithData$ ./util

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
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
	endl := []byte{}
	opts := NewOpts()
	cb := DefaultCallback(writer, opts)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, e := ParseXMLChan(datach, query, cb, opts); e != nil {
			t.Errorf("go routine failed with %v", e)
			return
		}
	}()
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

	wg.Wait()

	if writer.Len() < 1 {
		t.Errorf("no data written, check your source data and query")
	}
	fmt.Printf("%s\n", writer.String())
}

func TestParseXMLChan(t *testing.T) {
	xml := "<a>\n<b>hello\n</b>\n</a>"
	cdata := `
<root>
 <a> <![CDATA[Some important data too]]>
 </a>
 <b> <![CDATA[select *
              from table
			  where time <silly >;
    ]]> </b>
 <b> some other info</b>
</root>
`
	tests := []struct {
		xml     string
		xmlpath string
	}{
		{xml, "a>b"},
		{cdata, "root>a"},
		{cdata, "root>b"},
	}
	for _, test := range tests {
		parsexml(t, test.xml, test.xmlpath)
	}
}

func parsexml(t *testing.T, data string, xmlpath string) {
	datach := make(chan []byte)
	opts := NewOpts()
	opts.MakeOneLiner = false
	cb := DefaultCallback(os.Stdout, opts)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, e := ParseXMLChan(datach, xmlpath, cb, opts); e != nil {
			t.Errorf("go routine failed with %v", e)
		}
	}()
	datach <- []byte(data)
	//datach <- []byte{}
	close(datach)
	wg.Wait()
}
