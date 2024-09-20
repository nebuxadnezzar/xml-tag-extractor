package util

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type EVENT int

const MAXDOCSIZE = 32 * 1024 * 1024

const (
	PEEK EVENT = iota
	MID
	CDATA
	ENDTAG1
	ENDTAG2
	ENDDOC
	EOF
	NOOP = -1
)

type TagMap map[string]int

type ParserCallback func(string, EVENT) error

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

func (e EVENT) String() string {
	switch e {
	case PEEK:
		return "PEEK"
	case MID:
		return "MID"
	case CDATA:
		return "CDATA"
	case ENDTAG1:
		return "ENDTAG1"
	case ENDTAG2:
		return "ENDTAG2"
	case ENDDOC:
		return "ENDDOC"
	case EOF:
		return "EOF"
	case NOOP:
		return "NOOP"
	}
	return ``
}
