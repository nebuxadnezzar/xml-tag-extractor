package util

import (
	"fmt"
	"io"
	"os"
)

func MergeFiles(filenames []string, w io.Writer) (int64, error) {
	totalbytes := int64(0)
	var err error
	for _, name := range filenames {
		f, e := GetReader(name)
		if e != nil {
			err = fmt.Errorf("%w %w", err, e)
			continue
		}
		defer func(ff io.ReadCloser) { ff.Close() }(f)
		n, e := io.Copy(w, f)
		if e != nil {
			err = fmt.Errorf("%w %w", err, e)
			continue
		}
		totalbytes += n
	}
	return totalbytes, err
}

func Deletetemps(fa []*os.File) (err error) {
	if fa != nil {
		if len(fa) == 1 && fa[0] != os.Stdout {
			return nil
		}
		for i := len(fa); i > 0; i-- {
			f := fa[i-1]
			if f != nil {
				if e := os.Remove(f.Name()); e != nil {
					err = fmt.Errorf("%w %w", err, e)
				}
			}
		}
	}
	return
}

func Createtemps(sz int) (wa []*os.File, err error) {
	wa = make([]*os.File, sz)
	if sz == 1 {
		wa[0] = os.Stdout
		return wa, nil
	}
	for i := 0; i < sz; i++ {
		if f, e := os.CreateTemp(``, `wr*`); e == nil {
			wa[i] = f
		} else {
			err = e
			break
		}
	}
	if err != nil {
		Deletetemps(wa)
	}
	return wa, err
}

func GetReader(filename string) (io.ReadCloser, error) {
	if filename == os.Stdin.Name() {
		return os.Stdin, nil
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func CloseReader(reader io.ReadCloser, name string) error {
	if name != os.Stdin.Name() {
		return reader.Close()
	}
	return nil
}
