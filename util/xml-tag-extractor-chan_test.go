package util

// go test -v -count=1 -run ^TestParseXMLChan$ ./util
// go test -v -count=1 -run ^TestParseXMLChanWithData$ ./util

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"sync"
	"testing"
)

//go:embed data/consolidated.xml
var filedata []byte

func TestTagMapString(t *testing.T) {
	s := TagMapToStr(TagMap{"hi": 2})
	t.Logf("S: %s\n", s)
	if s == `` {
		t.Error("expected non-empty string")
	}
}

func TestEventToStr(t *testing.T) {

	if len(PEEK.String()) == 0 ||
		len(MID.String()) == 0 ||
		len(CDATA.String()) == 0 ||
		len(ENDTAG1.String()) == 0 ||
		len(ENDTAG2.String()) == 0 ||
		len(ENDDOC.String()) == 0 ||
		len(EOF.String()) == 0 {
		t.Error("expected non-empty string")
	}
}

func TestParseXMLChanWithData(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	reader := bufio.NewReader(bytes.NewBuffer(filedata))
	query := `CONSOLIDATED_LIST>INDIVIDUALS>INDIVIDUAL`
	datach := make(chan []byte)
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
			fmt.Printf("--> line# %d loop err: %v\n", cnt, err)
			break
		}
		datach <- []byte(s)
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
	buf := bytes.NewBuffer(nil)
	for _, test := range tests {
		parsexml(t, buf, test.xml, test.xmlpath)
	}
	println(fmt.Sprintf("\n%s\n", buf.String()))
}

func parsexml(t *testing.T, w io.Writer, data string, xmlpath string) {
	datach := make(chan []byte)
	opts := NewOpts()
	opts.MakeOneLiner = false
	cb := DefaultCallback(w, opts)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, e := ParseXMLChan(datach, xmlpath, cb, opts); e != nil {
			t.Errorf("go routine failed with %v", e)
		}
	}()
	datach <- []byte(data)
	close(datach)
	wg.Wait()
}
