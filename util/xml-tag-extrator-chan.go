package util

import (
	"bytes"
	"fmt"
)

func ParseXMLChan(datach chan []byte, errch chan error, query string, cb ParserCallback) (TagMap, error) {

	//defer close(errch)
	xb := bytes.NewBuffer(nil) //xml buffer
	st := NewPrefixstack()
	writeflag := false
	tagmap := map[string]int{}
	path := query + ">"
	cnt := 0
	for {
		buf, ok := <-datach
		if !ok {
			break
		}
		cnt++
		//fmt.Printf("%05d [%s]\n", cnt, string(buf))
		n := len(buf)
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
						//fmt.Printf("PFX %s %s %v\n", prefix, xb.String(), writeflag)
						if prefix == path && mr.event == ENDTAG2 {
							mr.event = ENDDOC
						}
					}

					//fmt.Printf("STACK: %s %v\n", st.String(), writeflag)
					if writeflag && cb != nil {
						if err := cb(xb.String(), mr.event); err != nil {
							errch <- err
							return nil, fmt.Errorf("error from XML parsing callback %w", err)
						}

						//if mr.event == MID || mr.event == ENDTAG1 {
						//	extractattr(xb.Bytes())
						//}
					}
					xb.Reset()
				default:
					xb.WriteByte(b)
					if xb.Len() > MAXDOCSIZE {
						err := fmt.Errorf("max. document size %d exceeded for path: %d %s", MAXDOCSIZE, xb.Len(), st.String())
						errch <- err
						return tagmap, err
					}
				}
			}
		}
	}
	return tagmap, nil
}
