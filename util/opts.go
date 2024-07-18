package util

import (
	"flag"
)

type Options struct {
	Files        []string
	XMLPaths     string
	MakeOneLiner bool
	AddRootTags  bool
	ShowHelp     bool
}

func NewOpts() *Options {
	return &Options{
		MakeOneLiner: true,
		AddRootTags:  false,
	}
}

func ParseArgs(args []string) *Options {
	opts := NewOpts()

	var xp string
	flag.StringVar(&xp, `xp`, ``, `CSV list of paths to tag(s) to extract, i.e. root:greeting OR root:greetings,root:story`)
	flag.StringVar(&xp, `xmlpaths`, ``, `CSV list of paths to tag(s) to extract, i.e. root:greeting OR root:greetings,root:story`)

	var ol, rt, hl bool
	flag.BoolVar(&ol, `ol`, true, `transform XML document into one-liner`)
	flag.BoolVar(&ol, `oneliner`, true, `transform XML document into one-liner`)
	flag.BoolVar(&rt, `rt`, false, `add root tags to each document to make correct XML document`)
	flag.BoolVar(&rt, `roottags`, false, `add root tags to each document to make correct XML document`)
	flag.BoolVar(&hl, `h`, false, `show help`)
	flag.BoolVar(&hl, `help`, false, `show help`)

	flag.CommandLine.Parse(args)

	opts.XMLPaths = xp
	opts.MakeOneLiner = ol
	opts.AddRootTags = rt
	opts.ShowHelp = hl
	opts.Files = flag.Args()

	return opts
}
