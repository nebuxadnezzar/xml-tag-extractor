package util

import (
	"flag"
)

type Options struct {
	Files        []string
	XMLPaths     string
	MakeOneLiner bool
	AddRootTags  bool
}

func newOpts() *Options {
	return &Options{
		MakeOneLiner: true,
		AddRootTags:  false,
	}
}

func ParseArgs(args []string) *Options {
	opts := newOpts()

	var xp string
	flag.StringVar(&xp, `xp`, ``, `CSV list of paths to tag(s) to extract, i.e. root:greeting OR root:greetings,root:story`)
	flag.StringVar(&xp, `xmlpaths`, ``, `CSV list of paths to tag(s) to extract, i.e. root:greeting OR root:greetings,root:story`)

	var ol, rt bool
	flag.BoolVar(&ol, `ol`, true, `transform XML document into one-liner`)
	flag.BoolVar(&ol, `oneliner`, true, `transform XML document into one-liner`)
	flag.BoolVar(&rt, `rt`, false, `add root tags to each document to make correct XML document`)
	flag.BoolVar(&rt, `roottags`, false, `add root tags to each document to make correct XML document`)

	flag.CommandLine.Parse(args)

	opts.XMLPaths = xp
	opts.MakeOneLiner = ol
	opts.AddRootTags = rt
	opts.Files = flag.Args()

	return opts
}
