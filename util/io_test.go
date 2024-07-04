package util

import (
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
