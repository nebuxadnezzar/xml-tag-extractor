package util

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestGetTemps(t *testing.T) {
	wa, err := Createtemps(2)
	if err != nil {
		t.Errorf("createtemps failed: %v", err)
		return
	}
	for i, w := range wa {
		t.Logf("%d %s\n", i, w.Name())
	}

	err = Deletetemps(wa)
	if err != nil {
		t.Logf("deletetemps failed: %v", err)
	}
}

func TestMergeTemps(t *testing.T) {
	wa, err := Createtemps(2)
	if err != nil {
		t.Errorf("%v", err)
	}
	defer func() {
		if wa != nil {
			Deletetemps(wa)
		}
	}()
	tmpnames := make([]string, 0, len(wa))
	for _, w := range wa {
		tmpnames = append(tmpnames, w.Name())
	}

	if _, err = MergeFiles(tmpnames, os.Stdout); err != nil {
		t.Errorf("%v", err)
	}
	fakereader := io.NopCloser(bytes.NewBuffer(nil))
	if _, err = GetReader(os.Stdin.Name()); err != nil {
		t.Errorf("%v", err)
	}
	if _, err = GetReader(`zz`); err == nil {
		t.Error("expected error")
	}
	CloseReader(fakereader, "fake")
	CloseReader(fakereader, os.Stdin.Name())
}
