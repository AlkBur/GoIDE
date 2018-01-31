package log

import (
	"testing"
)

func TestConsole(t *testing.T) {
	log1 := NewLogger(10000)
	log1.EnableFuncCallDepth(true)
	err := log1.SetLogger("console", "")
	if err != nil {
		t.Fatal(err)
	}
	testConsoleCalls(log1)

	log2 := NewLogger(100)
	log2.SetLogger("console", `{"level":3}`)
	testConsoleCalls(log2)
}

// Test console without color
func TestConsoleNoColor(t *testing.T) {
	log := NewLogger(100)
	err := log.SetLogger("console", `{"color":false}`)
	if err != nil {
		t.Fatal(err)
	}
	testConsoleCalls(log)
}

// Try each log level in decreasing order of priority.
func testConsoleCalls(bl *BaseLogger) {
	bl.Error("error")
	bl.Warn("warning")
	bl.Info("informational")
	bl.Debug("debug")
}