package utils

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func Exist(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

func ReadLineByFunc(file io.Reader, foo func(line []byte) error) error {
	if file == nil {
		return fmt.Errorf("ReadLines find reader is nil")
	}
	reader := bufio.NewReader(file)
	for {
		lines, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if err := foo(lines); err != nil {
			return err
		}
	}
	return nil
}

type FileInfo struct {
	File    string    `json:"file"`
	Size    int       `json:"size"`
	Time    time.Time `json:"time"`
	MD5     string    `json:"md5"`
	Changed bool      `json:"-"`
}

func ReadFileInfos(file string) ([]*FileInfo, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	result := make([]*FileInfo, 0)
	if err := json.NewDecoder(f).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func WriteFileInfos(files []*FileInfo, file string) error {
	sort.Slice(files, func(i, j int) bool {
		return files[i].File < files[j].File
	})
	output, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer output.Close()
	jsonOutput := json.NewEncoder(output)
	jsonOutput.SetEscapeHTML(false)
	jsonOutput.SetIndent("", "  ")
	return jsonOutput.Encode(files)
}

func SearchFile(files []*FileInfo, name string) *FileInfo {
	for _, v := range files {
		if v.File == name {
			return v
		}
	}
	return nil
}

func Md5File(file string) string {
	b, err := os.ReadFile(file)
	if err != nil {
		return ""
	}
	result := md5.Sum(b)
	return fmt.Sprintf("%x", result)
}

func MD5(file []byte) string {
	return fmt.Sprintf("%x", md5.Sum(file))
}

// WalkDirFiles 从路径dirPth下获取全部的文件.
func WalkDirFiles(dirPath string, filter func(base string, fileName string) error) error {
	var err error
	dirPath, err = filepath.Abs(dirPath)
	if err != nil {
		return err
	}
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info != nil && info.IsDir() {
			return nil
		}
		if err := filter(dirPath, path); err != nil {
			return err
		}
		return nil
	})
}

func ReadDir(dirname string) ([]fs.DirEntry, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	dirs, err := f.ReadDir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	return dirs, nil
}

func ReadDirs(dirname string, isAbs bool) ([]string, error) {
	dirs, err := ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, elem := range dirs {
		if elem.IsDir() {
			if isAbs {
				result = append(result, filepath.Join(dirname, elem.Name()))
			} else {
				result = append(result, elem.Name())
			}
		}
	}
	return result, nil
}

var _pwd = ""

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	_pwd = pwd
}
func FormatFile(filename string) string {
	if !filepath.IsAbs(filename) {
		return filename
	}
	rel, err := filepath.Rel(_pwd, filename)
	if err != nil {
		return filename
	}
	if strings.HasPrefix(rel, ".") {
		if rel == "." {
			return rel
		}
		return filename
	}
	return rel
}
