package util

import (
	"fmt"
	"regexp"
)

var (
	STAG = regexp.MustCompile(`(?is)^\s*<\s*([\w\-\.]+)[^>]*?([/]?>)\s*$`)
	ETAG = regexp.MustCompile(`(?is)^\s*</([\w\-\.]+)[^>]*?>\s*$`)
	ATTR = regexp.MustCompile(`(?is)\s+([\w\-]+)\s*=\s*"(.*?)"`)
)

type matchResult struct {
	tag     string
	event   EVENT
	matched bool
}

func matchtag(b []byte) matchResult {
	var mr matchResult
	//fmt.Printf("MATCHING: %s %v\n", string(b), STAG.Match(b))
	ma := STAG.FindSubmatch(b)
	l := len(ma)
	if l > 0 {
		mr.matched = true
		mr.tag = string(ma[1])
		mr.event = MID
		if l > 1 && ma[2][0] == '/' && ma[2][1] == '>' {
			mr.event = ENDTAG1
		}
		return mr
	}
	for i, v := range ma {
		fmt.Printf("\tT: %d %s\n", i, string(v))
	}
	return matchendtag(b)
}

func matchendtag(b []byte) matchResult {
	var mr matchResult
	//fmt.Printf("MATCHING: %s %v\n", string(b), ETAG.Match(b))
	ma := ETAG.FindSubmatch(b)
	l := len(ma)
	mr.matched = l > 0
	if mr.matched {
		mr.tag = string(ma[1])
		mr.event = ENDTAG2
		return mr
	}
	return mr
}

func extractattr(b []byte) map[string]string {
	fmt.Printf("\nEXTRACTING ATTR: %s\n", string(b))
	m := map[string]string{}
	ma := ATTR.FindSubmatch(b)
	if len(ma) > 0 {

	}
	for i, v := range ma {
		fmt.Printf("A: \t %d %s\n", i, string(v))
	}
	return m

}

func CreateOneLiner(s string) []byte {
	b := []byte(s)
	i := 0
	for j, k := 0, len(b); i < k && i+j < k; {
		ch := b[i]
		switch ch {
		case '\n', '\r':
			j++
			continue
		}
		b[i] = b[i+j]
		i++
	}
	return b[:i]
}
