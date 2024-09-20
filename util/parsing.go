package util

import (
	"bytes"
	"fmt"
	"regexp"
	"unicode"
)

var (
	STAG = regexp.MustCompile(`(?is)^\s*<([\w\-\.]+)[^>]*?([/]?>)\s*$`)
	ETAG = regexp.MustCompile(`(?is)^\s*</([\w\-\.]+)[^>]*?>\s*$`)
	ATTR = regexp.MustCompile(`(?is)\s+([\w\-]+)\s*=\s*"(.*?)"`)
	CDST = regexp.MustCompile(`(?is)^\s*<!\[CDATA\[(.*?)\s*$`)      // CDATA prefix
	CDAT = regexp.MustCompile(`(?is)^\s*<!\[CDATA\[(.*?)\]\]>\s*$`) // CDATA whole piece
	CDET = regexp.MustCompile(`(?is)^(.*?)\]\]>\s*$`)               // CDATA closing tag
)

type matchResult struct {
	tag       string
	event     EVENT
	matched   bool
	fullmatch bool
}

func matchcdata(b []byte) matchResult {
	var mr matchResult
	mr.event = NOOP
	ary := [3]*regexp.Regexp{CDAT, CDST, CDET}
	for i := range ary {
		ma := ary[i].FindSubmatch(b)
		if len(ma) > 0 {
			mr.matched = true
			mr.tag = string(ma[0])
			mr.event = CDATA
			mr.fullmatch = ary[i] == CDAT
			return mr
		}
	}
	return mr
}

func matchtag(b []byte) matchResult {
	var mr matchResult
	mr.event = NOOP
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
	return matchendtag(b)
}

func matchendtag(b []byte) matchResult {
	var mr matchResult
	mr.event = NOOP
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
	m := map[string]string{}
	ma := ATTR.FindAllSubmatch(b, -1)
	for _, v := range ma {
		if len(v) > 2 {
			m[string(v[1])] = string(v[2])
		}
	}
	return m
}

func maptoxml(m map[string]string) string {
	if m == nil || len(m) < 1 {
		return ``
	}
	b := bytes.NewBuffer(nil)
	for k, v := range m {
		b.WriteString(fmt.Sprintf(`<%s>%s</%s>`, k, v, k))
	}
	return b.String()
}

func CreateOneLiner(s string) []byte {
	b := []byte(s)
	i := 0
	for j, k := 0, len(b); i < k && i+j < k; {
		offset := i + j
		ch := b[offset]
		switch ch {
		case '\n', '\r':
			if offset+1 < k {
				ch := rune(b[offset+1])
				if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
					b[offset] = ' '
					goto SKIP
				}
			}
			j++
			continue
		}
	SKIP:
		b[i] = b[offset]
		i++
	}
	return b[:i]
}
