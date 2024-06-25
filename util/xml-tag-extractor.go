package util

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

type EVENT int

const MAXDOCSIZE = 32 * 1024 * 1024

const (
	PEEK EVENT = iota
	MID
	ENDTAG1
	ENDTAG2
	ENDDOC
)

type TagMap map[string]int

type ParserCallback func(string, EVENT) error

func ParseXML(reader io.Reader, query string, cb ParserCallback) (TagMap, error) {

	buf := make([]byte, 4096)
	xb := bytes.NewBuffer(nil) //xml buffer
	st := NewPrefixstack()
	writeflag := false
	tagmap := map[string]int{}
	path := query + ">"

	for {
		n, err := reader.Read(buf)
		//fmt.Fprintf(os.Stderr, "N: %d\n", n)
		if n > 0 {
			for i := 0; i < n; i++ {
				b := buf[i]
				switch b {
				case '<':
					if xb.Len() > 0 {
						//fmt.Printf("-> %s\n", xb.String())
						if writeflag && cb != nil {
							cb(xb.String(), PEEK)
						}
					}
					xb.Reset()
					xb.WriteByte(b)
				case '>':
					xb.WriteByte(b)
					mr := matchtag(xb.Bytes())
					//fmt.Printf("\n--> %s\n", xb.String())
					if mr.matched {
						prefix := st.String() + ">"
						switch mr.event {
						case MID:
							st.Push(mr.tag)
							prefix = st.String() + ">"
							writeflag = oktowrite(prefix, path, mr.tag, mr.event)
							updatemap(tagmap, prefix)
						case ENDTAG2:
							writeflag = oktowrite(prefix, path, mr.tag, mr.event)
							st.Pop()
						case ENDTAG1:
							writeflag = oktowrite(prefix, path, mr.tag, mr.event)
							//fmt.Fprintf(os.Stderr, "ok %v prefix: %s path: %s\n", writeflag, prefix, path)
							updatemap(tagmap, fmt.Sprintf("%s%s>", prefix, mr.tag))
						}

						if prefix == path && mr.event == ENDTAG2 {
							mr.event = ENDDOC
						}
					}

					//fmt.Printf("STACK: %s %v\n", st.String(), writeflag)
					if writeflag && cb != nil {
						cb(xb.String(), mr.event)
						//if mr.event == MID || mr.event == ENDTAG1 {
						//	extractattr(xb.Bytes())
						//}
					}
					xb.Reset()
				default:
					xb.WriteByte(b)
				}
			}
		}
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("%w", err)
			}
			break
		}
	}
	return tagmap, nil
}

func updatemap(m map[string]int, tag string) {
	if _, ok := m[tag]; ok {
		m[tag]++
	} else {
		m[tag] = 1
	}
}

func oktowrite(prefix, path, tag string, evt EVENT) bool {
	if path == "" || prefix == "" {
		return false
	}
	pfx := prefix
	if evt == ENDTAG1 {
		pfx = fmt.Sprintf("%s%s>", prefix, tag)
	}
	//fmt.Printf("\nEVENT: %d PFX: %s PATH: %s\n", evt, pfx, path)
	return strings.HasPrefix(pfx, path)
}

func TagMapToStr(m TagMap) string {

	kk := make([]string, 0, len(m))
	for k := range m {
		kk = append(kk, k)
	}
	sort.Strings(kk)
	buf := bytes.NewBuffer(nil)
	for _, k := range kk {
		buf.WriteString(fmt.Sprintf("%s\t\t%d\n", strings.Replace(strings.TrimRight(k, ">"), ">", ":", -1), m[k]))
	}
	return buf.String()
}
