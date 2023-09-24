package utils

import (
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"

	"github.com/anthony-dong/mac_tools/logs"
)

type Md5Writer struct {
	md5_ hash.Hash
	w    io.Writer
	r    io.Reader
}

func NewMd5Writer(w io.Writer) *Md5Writer {
	r := Md5Writer{}
	r.md5_ = md5.New()
	r.w = w
	return &r
}

func NewMd5Reader(r_ io.Reader) *Md5Writer {
	r := Md5Writer{}
	r.md5_ = md5.New()
	r.r = r_
	return &r
}

func (r *Md5Writer) Sum() string {
	vv := r.md5_.Sum(nil)
	return fmt.Sprintf("%x", vv)
}

func (r *Md5Writer) Write(p []byte) (n int, err error) {
	n, err = r.w.Write(p)
	if n > 0 {
		_, _ = r.md5_.Write(p[:n])
	}
	return
}

func (r *Md5Writer) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if n > 0 {
		_, _ = r.md5_.Write(p[:n])
	}
	return
}

func CopyFile(from string, to string, fileinfo *FileInfo) error {
	fromFileStat, err := os.Stat(from)
	if err != nil {
		logs.Error("open file find err: %v, file: %s", err, FormatFile(from))
		return err
	}
	if fromFileStat.ModTime().Unix() == fileinfo.Time.Unix() {
		return nil
	}
	fileinfo.Changed = true
	fromFile, err := os.Open(from)
	if err != nil {
		logs.Error("open file find err: %v, file: %s", err, FormatFile(from))
		return err
	}
	defer fromFile.Close()
	if dir := filepath.Dir(to); !Exist(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logs.Error("create dir find err: %v, dir: %s", err, FormatFile(dir))
			return err
		}
	}
	toFile, err := os.Create(to)
	if err != nil {
		logs.Error("create file find err: %v, file: %s", err, FormatFile(to))
		return err
	}
	defer toFile.Close()
	md5_ := NewMd5Reader(fromFile)
	size, err := io.Copy(toFile, md5_)
	if err != nil {
		logs.Error("copy file find err: %v, file: %s -> %s", err, FormatFile(from), FormatFile(to))
		return err
	}
	oldtime := fileinfo.Time.Unix()
	oldMd5 := fileinfo.MD5
	fileinfo.Time = fromFileStat.ModTime()
	fileinfo.Size = int(size)
	fileinfo.MD5 = md5_.Sum()
	if oldMd5 != "" {
		logs.Info("copy file %s -> %s (%d->%d %s->%s)", FormatFile(from), FormatFile(to), oldtime, fileinfo.Time.Unix(), oldMd5, fileinfo.MD5)
	} else {
		logs.Info("copy file %s -> %s (%d-%s)", FormatFile(from), FormatFile(to), fileinfo.Time.Unix(), fileinfo.MD5)
	}
	return nil
}
