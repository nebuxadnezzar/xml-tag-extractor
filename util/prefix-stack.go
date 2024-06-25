package util

import (
	"bytes"
	"fmt"
)

type prefixstack struct {
	stack []byte
	last  int
	top   int
}

func NewPrefixstack() *prefixstack {
	return &prefixstack{
		stack: make([]byte, 32000),
	}
}
func (p *prefixstack) Top() (string, error) {
	b := bytes.NewBuffer(nil)
	if p.last == 0 {
		return "", fmt.Errorf("stack is empty")
	}
	//fmt.Printf("TOP %d LAST %d\n", p.top, p.last)
	for i := p.top; i < p.last; i++ {
		b.WriteByte(p.stack[i])
	}
	return b.String(), nil
}

func (p *prefixstack) Pop() error {
	if p.last == 0 {
		return fmt.Errorf("stack is empty")
	}
	// popping the last item
	if p.top == 0 && p.last > 0 {
		p.last = 0
		return nil
	}
	for ; p.top < p.last; p.last-- {
	}
	p.last -= 1
	p.top -= 2
	for ; p.top > -1 && p.stack[p.top] != '>'; p.top-- {
	}
	if p.top < 0 {
		p.top = 0
	}
	if p.stack[p.top] == '>' {
		p.top++
	}
	return nil
}

func (p *prefixstack) Push(s string) error {
	l := len(s)
	if !(l+p.last+1 < len(p.stack)) {
		return fmt.Errorf("%s too big: only %d bytes left on the stack", s, len(p.stack)-p.last)
	}
	n := 0
	b := []byte(s)
	p.top = p.last
	if p.top > 0 {
		p.top = p.last + 1
		p.stack[p.last] = '>'
		p.last++
	}
	for i, k := 0, len(p.stack); i+p.last < k && i < l; i++ {
		p.stack[i+p.last] = b[i]
		n++
	}
	p.last += n
	//fmt.Printf("after push: %s top %d [%c] last %d [%c]\n", s, p.top, p.stack[p.top], p.last, p.stack[p.last])
	return nil
}

func (p *prefixstack) String() string {
	return string(p.stack[:p.last])
}
