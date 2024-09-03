package util

import (
	"bytes"
	"fmt"
)

func ParseXMLChan(datach chan []byte, query string, cb ParserCallback, opts *Options) (TagMap, error) {

	xb := bytes.NewBuffer(nil) //xml buffer
	st := NewPrefixstack()
	writeflag := false
	tagmap := map[string]int{}
	path := query + ">"
	cnt := 0
	for {
		buf, ok := <-datach
		if !ok {
			cb(``, EOF)
			break
		}
		cnt++
		//fmt.Printf("%05d [%s]\n", cnt, string(buf))
		n := len(buf)
		if n > 0 {
			for i := 0; i < n; i++ {
				b := buf[i]
				switch b {
				case '\n', '\r':
					if opts.MakeOneLiner {
						xb.WriteByte(0x20)
					} else {
						xb.WriteByte(b)
					}
				case '<':
					if xb.Len() > 0 {
						mr := matchcdata(xb.Bytes())
						if mr.matched {
							//println(fmt.Sprintf("M1 CDATA EVT: %s %s", mr.event.String(), mr.tag))
							xb.WriteByte(b)
							continue
						}
						//fmt.Printf("-> %s\n", xb.String())
						//if bytes.Contains(xb.Bytes(), []byte(`<![CDATA[`)) {
						//fmt.Printf("CDATA OUTSIDE OF CDATA EVENT: %s\n", xb.String())
						//continue
						//}
						if writeflag && cb != nil {
							cb(xb.String(), PEEK)
						}
					}
					xb.Reset()
					xb.WriteByte(b)
				case '>':
					xb.WriteByte(b)
					mr := matchtag(xb.Bytes())
					//fmt.Printf("\n--> MR: %#v %s %v %s\n", mr, mr.event.String(), writeflag, xb.String())
					if mr.matched {
						prefix := st.String() + ">"
						switch mr.event {
						case CDATA:
							writeflag = oktowrite(prefix, path, mr.tag, mr.event)
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
					} else {
						mr = matchcdata(xb.Bytes())
						if mr.matched {
							//println(fmt.Sprintf("M2 CDATA EVT: %s %s %v", mr.event.String(), mr.tag, mr.fullmatch))
							if !mr.fullmatch {
								continue
							}
						}
					}

					//fmt.Printf("STACK: %s %s %v %s\n", st.String(), mr.event.String(), writeflag, xb.String())
					if writeflag && cb != nil {
						if !opts.AttributesToElements || mr.event == ENDDOC {
							if err := cb(xb.String(), mr.event); err != nil {
								return nil, fmt.Errorf("error from XML parsing callback %w", err)
							}
						} else if xatr := maptoxml(extractattr(xb.Bytes())); len(xatr) > 0 {
							endtag := ``
							if mr.event == ENDTAG1 {
								endtag = fmt.Sprintf(`</%s>`, mr.tag)
							}
							if err := cb(fmt.Sprintf(`<%s>%s%s`, mr.tag, xatr, endtag), mr.event); err != nil {
								return nil, fmt.Errorf("error from XML parsing callback %w", err)
							}
						} else {
							if err := cb(xb.String(), mr.event); err != nil {
								return nil, fmt.Errorf("error from XML parsing callback %w", err)
							}
						}
					}
					xb.Reset()
				default:
					xb.WriteByte(b)
					if xb.Len() > MAXDOCSIZE {
						err := fmt.Errorf("max. document size %d exceeded for path: %d %s", MAXDOCSIZE, xb.Len(), st.String())
						//errch <- err
						return tagmap, err
					}
				}
			}
		}
	}
	return tagmap, nil
}
