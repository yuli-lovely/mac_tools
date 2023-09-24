package logs

import "testing"

func TestDebug(t *testing.T) {
	t.Log("debug")
	SelLevel(LevelDebug)
	Debug("1")
	Info("2")
	Error("3")

	t.Log("info")
	SelLevel(LevelInfo)
	Debug("1")
	Info("2")
	Error("3")

	t.Log("error")
	SelLevel(LevelError)
	Debug("1")
	Info("2")
	Error("3")
}
