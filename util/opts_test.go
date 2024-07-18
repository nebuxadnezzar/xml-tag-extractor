package util

import (
	"fmt"
	"testing"
)

func TestParseOptions(t *testing.T) {
	args := []string{"-xp=root:hello", "-ol", "-rt", "root", "my.xml"}
	opts := ParseArgs(args)
	fmt.Printf("%#v\n", opts)
	if !(opts.Files[0] == "my.xml" && opts.XMLPaths == "root:hello" && opts.RootTagName == "root" && opts.MakeOneLiner) {
		t.Errorf("one of expected values is missing: %s %s %s %v", opts.Files[0], opts.XMLPaths, opts.RootTagName, opts.MakeOneLiner)
	}
}
