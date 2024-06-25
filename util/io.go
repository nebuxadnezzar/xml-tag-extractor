package util

import (
	"io"
	"os"
)

func deletetemps(fa []*os.File) {
	if fa != nil {
		for i := len(fa); i > 0; i-- {
			f := fa[i]
			if f != nil {
				os.Remove(f.Name())
			}
		}
	}
}

func createtemps(sz int) (wa []*os.File, err error) {
	wa = make([]*os.File, sz)
	for i := 0; i < sz; i++ {
		if f, e := os.CreateTemp(``, `wr*`); e == nil {
			wa[i] = f
		} else {
			err = e
			break
		}
	}
	if err != nil {
		deletetemps(wa)
	}
	return wa, err
}

func GetReader(filename string) (io.ReadCloser, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return f, nil
}
