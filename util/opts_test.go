package util

import (
	"fmt"
	"testing"
)

func TestParseOptions(t *testing.T) {
	args := []string{"-xp=root:hello", "-ol", "-rt", "my.xml"}
	opts := ParseArgs(args)
	fmt.Printf("%#v\n", opts)
	if !(opts.Files[0] == "my.xml" && opts.XMLPaths == "root:hello" && opts.AddRootTags && opts.MakeOneLiner) {
		t.Errorf("one of expected values is missing: %s %s %v %v", opts.Files[0], opts.XMLPaths, opts.AddRootTags, opts.MakeOneLiner)
	}
}
