package logs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var debug = false
var _level = LevelInfo

func init() {
	if os.Getenv("DEBUG") == "1" {
		SelLevel(LevelDebug)
	}
}

func SelLevel(level int) {
	if level == LevelDebug {
		debug = true
	} else {
		debug = false
	}
	_level = level
}

const LevelDebug = 0
const LevelInfo = 1
const LevelError = 2

func Debug(format string, v ...interface{}) {
	log(LevelDebug, format, v...)
}

func Info(format string, v ...interface{}) {
	log(LevelInfo, format, v...)
}

func Error(format string, v ...interface{}) {
	log(LevelError, format, v...)
}

func log(level int, format string, v ...interface{}) {
	if level < _level {
		return
	}
	now := time.Now().Format("15:04:05.000")
	out := strings.Builder{}
	switch level {
	case LevelDebug:
		out.WriteString("[DEBUG] ")
	case LevelInfo:
		out.WriteString("[INFO] ")
	case LevelError:
		out.WriteString("[ERROR] ")
	default:
		return
	}
	out.WriteString(now)
	out.WriteString(" ")
	if debug {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
		out.WriteString(filepath.Base(file))
		out.WriteString(":")
		out.WriteString(strconv.FormatInt(int64(line), 10))
		out.WriteString(" ")
	}
	out.WriteString(fmt.Sprintf(format, v...))
	fmt.Println(out.String())
}
