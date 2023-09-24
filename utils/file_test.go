package utils

import (
	"testing"

	"github.com/anthony-dong/mac_tools/logs"
	"github.com/stretchr/testify/assert"
)

func TestExist(t *testing.T) {
	t.Log(Exist("/Users/yuli/go/src/github.com/anthony-dong/mac_tools"))
}

func Test_Md5File(t *testing.T) {
	// 95aef247c291d56d20243ca28a38cd7e
	assert.Equal(t, Md5File("/Users/yuli/go/src/github.com/anthony-dong/mac_tools/main.go"), "95aef247c291d56d20243ca28a38cd7e")
	t.Log(Md5File("/Users/yuli/go/src/github.com/anthony-dong/mac_tools/main.go"))
}

// ReadDir

func Test_ReadDir(t *testing.T) {
	dirs, err := ReadDirs("/Users/yuli/go/src/github.com/anthony-dong/mac_tools", false)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range dirs {
		logs.Info("file: %s", v)
	}
}
