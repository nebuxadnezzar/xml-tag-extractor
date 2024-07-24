//go:build !perf

package running

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/nebuxadnezzar/xml-tag-extractor/util"
)

func Runit() {
	opts := util.NewOpts()
	if len(os.Args) > 1 {
		opts = util.ParseArgs(os.Args[1:])
	}
	fmt.Fprintf(os.Stderr, "OPTS: %#v\n", opts)

	if opts.ShowHelp {
		showUsageAndExit()
	}
	path := ""
	if opts.XMLPaths != `` {
		path = strings.Replace(opts.XMLPaths, ":", ">", -1)
	}
	fmt.Fprintf(os.Stderr, "PATH: %s\n", path)
	pp := strings.Split(path, ",")

	filename := os.Stdin.Name()
	if len(opts.Files) > 0 {
		filename = opts.Files[0]
	}
	if len(pp) > 1 {
		os.Exit(run1(filename, pp, opts))
	} else {
		os.Exit(run(filename, path, opts))
	}
}

func run(filename, path string, opts *util.Options) (status int) {
	reader, err := util.GetReader(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening: %s %v\n", os.Args[1], err)
		return 2
	}
	defer util.CloseReader(reader, filename)

	printHeaderOrFooter(os.Stdout, filepath.Base(filename), opts, true)

	tagmap, err := util.ParseXML(reader, path, util.DefaultCallback(os.Stdout, opts), opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating tagmap: %v\n", err)
		return 3
	}
	printHeaderOrFooter(os.Stdout, filename, opts, false)
	if path == `` {
		fmt.Printf("%s", util.TagMapToStr(tagmap))
	}
	return 0

}

func run1(filename string, pp []string, opts *util.Options) (status int) {
	reader, err := util.GetReader(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening: %s %v\n", os.Args[1], err)
		return 2
	}
	defer util.CloseReader(reader, filename)

	wa, err := util.Createtemps(len(pp))
	defer func() {
		if wa != nil {
			util.Deletetemps(wa)
		}
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "createtemps failed: %v", err)
		return 1
	}

	rd := bufio.NewReader(reader)
	wg := new(sync.WaitGroup)
	errch := make(chan error)
	datachs := make([]chan []byte, len(pp))
	for i := range datachs {
		datachs[i] = make(chan []byte, 16)
	}

	for i, p := range pp {
		wg.Add(1)
		go func(path string, ii int, datach chan []byte, o *util.Options) {
			defer wg.Done()
			//println("subpath", path)
			cb := util.DefaultCallback(wa[ii], o)

			if tagmap, err := util.ParseXMLChan(datach, errch, path, cb, opts); err == nil {
				if path == `` {
					fmt.Println(util.TagMapToStr(tagmap))
				}
			}
		}(p, i, datachs[i], opts)
	}
	go func() {
		for e := range errch {
			fmt.Fprintf(os.Stderr, ">> %v\n", e)
		}
		wg.Wait()
		close(errch)
	}()

	cnt := 0
	for {
		s, err := rd.ReadString('\n')
		cnt++
		if len(s) > 0 {
			//fmt.Printf("%05d %d Sending [%s]", cnt, len(datachs), s)
			for i, k := 0, len(datachs); i < k; i++ {
				ch := datachs[i]
				//print(cnt, " channeling", s)
				ch <- []byte(s)
				ch <- []byte{}
			}
		}
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "line# %d loop err: %v\n", cnt, err)
			}
			break
		}
	}

	for _, ch := range datachs {
		close(ch)
	}

	tmpnames := make([]string, 0, len(pp))
	for _, w := range wa {
		tmpnames = append(tmpnames, w.Name())
	}
	printHeaderOrFooter(os.Stdout, filepath.Base(filename), opts, true)
	total, err := util.MergeFiles(tmpnames, os.Stdout)
	printHeaderOrFooter(os.Stdout, filepath.Base(filename), opts, false)

	fmt.Fprintf(os.Stderr, "\ntotal bytes: %d\n", total)
	if err != nil {
		fmt.Fprintf(os.Stderr, "copy error: %v\n", err)
	}
	return 0
}

func showUsageAndExit() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] XML-file-path\nOPTIONS:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(0)
}

func printHeaderOrFooter(w io.Writer, source string, opts *util.Options, isheader bool) {
	if opts.RootTagName != `` && opts.XMLPaths != `` {
		if isheader {
			w.Write([]byte(fmt.Sprintf("<%s src=\"%s\" ts=\"%s\">\n",
				opts.RootTagName,
				source,
				time.Now().Local().Format(`2006-01-02 15:04:05`))))
		} else {
			w.Write([]byte(fmt.Sprintf("\n</%s>", opts.RootTagName)))
		}
	}
}

/*

go run cmd/main.go ~/test-data/consolidated.xml CONSOLIDATED_LIST:INDIVIDUALS:INDIVIDUAL,CONSOLIDATED_LIST:ENTITIES:ENTITY >

*/
