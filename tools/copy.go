package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/anthony-dong/mac_tools/logs"
	"github.com/anthony-dong/mac_tools/utils"
	"github.com/spf13/cobra"
)

func NewCopyCmd() (*cobra.Command, error) {
	helper := CopyHelper{
		done: make(chan bool),
	}
	cmd := &cobra.Command{
		Use: "copy",
		RunE: func(cmd *cobra.Command, args []string) error {
			RegisteDeferTask(func() {
				helper.Close()
			})
			return helper.Run()
		},
	}
	cmd.Flags().StringVarP(&helper.From, "from", "f", "", "the from dir")
	cmd.Flags().StringVarP(&helper.From, "to", "f", "", "the to dir")
	return cmd, nil
}

var ErrorDone = fmt.Errorf(`ctx is done`)

type CopyHelper struct {
	From     string
	To       string
	done     chan bool
	doneOnce sync.Once
}

type MetaInfo struct {
	FileName  string
	FileInfos []*utils.FileInfo
}

func NewCopyHelper(from string, to string) *CopyHelper {
	return &CopyHelper{
		From: from,
		To:   to,
		done: make(chan bool),
	}
}

func (c *CopyHelper) assert() error {
	assertFile := filepath.Join(c.To, "mac_tools.assert")
	want := fmt.Sprintf("%s -> %s", c.From, c.To)
	if !utils.Exist(assertFile) {
		return os.WriteFile(assertFile, []byte(want), 0644)
	}
	got, err := os.ReadFile(assertFile)
	if err != nil {
		return err
	}
	if string(got) == want {
		return nil
	}
	return fmt.Errorf(`must copy %s`, got)
}

func (c *CopyHelper) Close() {
	c.doneOnce.Do(func() {
		close(c.done)
	})
}

func (c *CopyHelper) init() error {
	var err error
	if c.From, err = filepath.Abs(c.From); err != nil {
		return err
	}
	if c.To, err = filepath.Abs(c.To); err != nil {
		return err
	}
	logs.Info("copy dir %s -> %s", utils.FormatFile(c.From), utils.FormatFile(c.To))
	if !utils.Exist(c.From) {
		return fmt.Errorf(`not exist file %s`, c.From)
	}
	if !utils.Exist(c.To) {
		return fmt.Errorf(`not exist file %s`, c.To)
	}
	if err := c.assert(); err != nil {
		return err
	}
	return nil
}

func (c *CopyHelper) Run() error {
	if err := c.init(); err != nil {
		return err
	}
	froms, err := utils.ReadDirs(c.From, false)
	if err != nil {
		return nil
	}
	for _, from := range froms {
		fromDirAbs := filepath.Join(c.From, from)
		toDirAbs := filepath.Join(c.To, from)
		if strings.HasPrefix(fromDirAbs, c.To) {
			continue
		}
		mateinfo := &MetaInfo{}
		mateinfo.FileName = filepath.Join(toDirAbs, ".metainfo.json")
		mateinfo.FileInfos = make([]*utils.FileInfo, 0)
		if utils.Exist(mateinfo.FileName) {
			if m, err := utils.ReadFileInfos(mateinfo.FileName); err != nil {
				return err
			} else {
				logs.Debug("read meta info file %s, len %d", utils.FormatFile(mateinfo.FileName), len(m))
				mateinfo.FileInfos = m
			}
		}
		logs.Info("copy dir %s -> %s", utils.FormatFile(fromDirAbs), utils.FormatFile(toDirAbs))
		isWrite := false
		err := utils.WalkDirFiles(fromDirAbs, func(base, fromFile string) error {
			select {
			case <-c.done:
				return ErrorDone
			default:
			}
			relName, err := filepath.Rel(base, fromFile)
			if err != nil {
				return err
			}
			isCreate := false
			fileInfo := utils.SearchFile(mateinfo.FileInfos, relName)
			if fileInfo == nil {
				isCreate = true
				fileInfo = &utils.FileInfo{File: relName}
			}
			fileInfo.Changed = false
			toFile := filepath.Join(toDirAbs, relName)
			if err := utils.CopyFile(fromFile, toFile, fileInfo); err != nil {
				return err
			}
			if !isCreate {
				if fileInfo.Changed {
					isWrite = true
				}
				return nil
			}
			isWrite = true
			mateinfo.FileInfos = append(mateinfo.FileInfos, fileInfo)
			return nil
		})
		doWrite := func() error {
			if isWrite && len(mateinfo.FileInfos) > 0 {
				if err := utils.WriteFileInfos(mateinfo.FileInfos, mateinfo.FileName); err != nil {
					logs.Error("write file %s find err %v", utils.FormatFile(mateinfo.FileName), err)
					return err
				}
				logs.Info("write meta info file %s, len %d", utils.FormatFile(mateinfo.FileName), len(mateinfo.FileInfos))
			}
			return nil
		}
		if err != nil {
			if err == ErrorDone {
				return doWrite()
			}
			logs.Error("copy dir %s -> %s find err: %v", utils.FormatFile(fromDirAbs), utils.FormatFile(toDirAbs), err)
			return err
		}
		if err := doWrite(); err != nil {
			return err
		} else {
			continue
		}
	}
	return nil
}
