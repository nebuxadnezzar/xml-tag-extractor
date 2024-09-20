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
	//fmt.Fprintf(os.Stderr, "OPTS: %#v\n", opts)

	if opts.ShowHelp {
		showUsageAndExit()
	}
	path := ""
	if opts.XMLPaths != `` {
		path = strings.Replace(opts.XMLPaths, ":", ">", -1)
	}
	pp := strings.Split(path, ",")

	filename := os.Stdin.Name()
	if len(opts.Files) > 0 {
		filename = opts.Files[0]
	}
	os.Exit(run1(filename, pp, opts))
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

	chnum := 1
	if opts.Boost {
		chnum = 24
	}

	rd := bufio.NewReader(reader)
	wg := new(sync.WaitGroup)
	datachs := make([]chan []byte, len(pp))
	for i := range datachs {
		datachs[i] = make(chan []byte, chnum)
	}

	for i, p := range pp {
		wg.Add(1)
		go func(path string, ii int, datach chan []byte, o *util.Options) {
			defer wg.Done()
			cb := util.DefaultCallback(wa[ii], o)

			if tagmap, err := util.ParseXMLChan(datach, path, cb, opts); err == nil {
				if path == `` {
					fmt.Printf("\n%s\n", util.TagMapToStr(tagmap))
				}
			} else {
				fmt.Fprintf(os.Stderr, "routine %d returned error %v\n", ii, err)
			}
		}(p, i, datachs[i], opts)
	}

	cnt := 0
	totalbytes := 0
	xmlPathNotPresent := opts.XMLPaths == ``
	for {
		s, err := rd.ReadString('>')
		cnt++
		if len(s) > 0 {
			if xmlPathNotPresent {
				totalbytes += len(s)
				fmt.Fprintf(os.Stderr, "\rbytes read: %d", totalbytes)
			}
			for i, k := 0, len(datachs); i < k; i++ {
				datachs[i] <- []byte(s)
			}
		}
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "line# %d loop err: %v\n", cnt, err)
			}
			break
		}
	}

	// give a pause to drain buffered channels
	time.Sleep(100 * time.Millisecond)

	for i, ch := range datachs {
		close(ch)
		for range ch {
			d := <-ch
			if len(d) > 0 {
				fmt.Fprintf(os.Stderr, "channel %d still has buffer: [%s]\n", i, string(d))
			}
		}
	}
	wg.Wait()

	tmpnames := make([]string, 0, len(pp))
	for _, w := range wa {
		tmpnames = append(tmpnames, w.Name())
	}
	printHeaderOrFooter(os.Stdout, filepath.Base(filename), opts, true)
	total, err := util.MergeFiles(tmpnames, os.Stdout)
	printHeaderOrFooter(os.Stdout, filepath.Base(filename), opts, false)

	fmt.Fprintf(os.Stderr, "\nbytes merged: %d\n", total)
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
