package tools

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/anthony-dong/mac_tools/logs"
	"github.com/spf13/cobra"
)

var _signal chan os.Signal

var deferTask []func()

func InitDeferTask() {
	_signal = make(chan os.Signal, 1)
	signal.Notify(_signal, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
}

func RegisteDeferTask(task func()) {
	if task == nil {
		panic("defer task is nil")
	}
	deferTask = append(deferTask, task)
}

func RunDeferTask() {
	cc := <-_signal
	if cc != nil {
		logs.Info("watch signal: %s", cc)
	}
	for _, task := range deferTask {
		task()
	}
}

func AddCmd(cmd *cobra.Command, foo func() (*cobra.Command, error)) error {
	subCmd, err := foo()
	if err != nil {
		return err
	}
	cmd.AddCommand(subCmd)
	return nil
}
