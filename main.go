package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/anthony-dong/mac_tools/logs"
	"github.com/anthony-dong/mac_tools/tools"
	"github.com/spf13/cobra"
)

func run() {
	cmd := cobra.Command{
		Use:     "mac_tools",
		Version: "1.0.0",
	}
	err := cmd.ExecuteContext(context.Background())
	if err != nil {
		panic(err)
	}
}

func main() {
	from := ""
	to := ""
	flag.StringVar(&from, "from", "", "@dir: 表示当前目录")
	flag.StringVar(&to, "to", "", "复制的目录")

	flag.Parse()
	if from == "" || to == "" {
		flag.Usage()
		panic("from、to is empty")
	}
	if from == "@cur" {
		if path, err := os.Getwd(); err != nil {
			panic(err)
		} else {
			from = path
		}
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	copy := tools.NewCopyHelper(from, to)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		cc := <-c
		if cc != nil {
			logs.Info("watch signal: %s", cc)
		}
		copy.Close()
	}()
	go func() {
		defer wg.Done()
		if err := copy.Run(); err != nil {
			panic(err)
		} else {
			close(c)
		}
	}()
	wg.Wait()
}
