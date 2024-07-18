package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/nebuxadnezzar/xml-tag-extractor/util"
)

func main() {

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
		os.Exit(run1(filename, pp))
	} else {
		os.Exit(run(filename, path))
	}
}

func run(filename, path string) (status int) {
	reader, err := util.GetReader(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening: %s %v\n", os.Args[1], err)
		return 2
	}
	defer util.CloseReader(reader, filename)

	tagmap, err := util.ParseXML(reader, path, util.DefaultCallback(os.Stdout))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating tagmap: %v\n", err)
		return 3
	}
	if path == `` {
		fmt.Printf("%s", util.TagMapToStr(tagmap))
	}
	return 0

}

func run1(filename string, pp []string) (status int) {
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
		go func(path string, ii int, datach chan []byte) {
			defer wg.Done()
			//println("subpath", path)
			cb := util.DefaultCallback(wa[ii])

			if tagmap, err := util.ParseXMLChan(datach, errch, path, cb); err == nil {
				if path == `` {
					fmt.Println(util.TagMapToStr(tagmap))
				}
			}
		}(p, i, datachs[i])
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
	total, err := util.MergeFiles(tmpnames, os.Stdout)
	fmt.Fprintf(os.Stderr, "\ntotal bytes: %d\n", total)
	if err != nil {
		fmt.Fprintf(os.Stderr, "copy error: %v\n", err)
	}
	return 0
}

func showUsageAndExit() {
	fmt.Printf("Usage: %s file start-tag", os.Args[0])
	os.Exit(0)
}

/*

go run cmd/main.go ~/test-data/consolidated.xml CONSOLIDATED_LIST:INDIVIDUALS:INDIVIDUAL,CONSOLIDATED_LIST:ENTITIES:ENTITY >

*/
