package util

import (
	"flag"
)

type Options struct {
	// XML file name. comes as the rest afte all options are processed
	Files []string

	// XMLPaths are CSV string. Example: root:greetings OR root:greetings:hello
	// where root, greetings and hello are tags to be parsed out
	XMLPaths string

	// convert parsed out multiline XML document to one liner
	MakeOneLiner bool

	// convert attributes to elements
	AttributesToElements bool

	// add root tag to resulting XML records
	RootTagName string

	// show help
	ShowHelp bool
}

func NewOpts() *Options {
	return &Options{
		MakeOneLiner: true,
	}
}

func ParseArgs(args []string) *Options {
	opts := NewOpts()

	var xp, rt string
	flag.StringVar(&xp, `xp`, ``, `CSV list of paths to tag(s) to extract, i.e. root:greeting OR root:greetings,root:story`)
	flag.StringVar(&xp, `xml.paths`, ``, `CSV list of paths to tag(s) to extract, i.e. root:greeting OR root:greetings,root:story`)
	flag.StringVar(&rt, `rt`, ``, `add provided root tags to each document to make correct XML document`)
	flag.StringVar(&rt, `root.tag`, ``, `add provided root tags to each document to make correct XML document`)

	var ol, hl, ca bool
	flag.BoolVar(&ca, `ca`, false, `convert attributes to elements`)
	flag.BoolVar(&ca, `convert.attributes`, false, `convert attributes to elements`)
	flag.BoolVar(&ol, `ol`, false, `transform XML document into one-liner`)
	flag.BoolVar(&ol, `one.liner`, false, `transform XML document into one-liner`)
	flag.BoolVar(&hl, `h`, false, `show help`)
	flag.BoolVar(&hl, `help`, false, `show help`)

	flag.CommandLine.Parse(args)

	opts.XMLPaths = xp
	opts.AttributesToElements = ca
	opts.MakeOneLiner = ol
	opts.RootTagName = rt
	opts.ShowHelp = hl
	opts.Files = flag.Args()

	return opts
}
