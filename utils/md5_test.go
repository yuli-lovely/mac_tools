package utils

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNewMd5Writer(t *testing.T) {
	file, err := os.Open("/Users/yuli/go/src/github.com/anthony-dong/mac_tools/main.go")
	if err != nil {
		t.Fatal(err)
	}

	mw := NewMd5Reader(file)

	_, err = ioutil.ReadAll(mw)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(mw.Sum())
	// 676040c722f4e2d50137e3206b34f122
}
