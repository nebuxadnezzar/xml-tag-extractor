package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/nebuxadnezzar/xml-tag-extractor/util"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s file start-tag", os.Args[0])
		os.Exit(1)
	}

	path := ""
	if len(os.Args) > 2 {
		path = strings.Replace(os.Args[2], ":", ">", -1)
	}
	fmt.Fprintf(os.Stderr, "PATH: %s\n", path)
	pp := strings.Split(path, ",")

	wg := new(sync.WaitGroup)
	for _, p := range pp {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			reader, err := util.GetReader(os.Args[1])
			if err != nil {
				fmt.Printf("error opening: %s %v\n", os.Args[1], err)
				os.Exit(2)
			}
			defer reader.Close()
			println("subpath", path)
			if tagmap, err := util.ParseXML(reader, path, util.DefaultCallback(os.Stdout)); err == nil {
				if path == `` {
					fmt.Printf(util.TagMapToStr(tagmap))
				}
			} else {
				fmt.Fprintf(os.Stderr, "err: %v\n", err)
			}
		}(p)
	}
	wg.Wait()
}
