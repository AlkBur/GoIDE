package log

import (
	"encoding/json"
	"os"
	"time"
)

var colors = []brush{
	newBrush("1;31"), // Error              red
	newBrush("1;33"), // Warning            yellow
	newBrush("1;34"), // Informational      blue
	newBrush("1;37"), // Debug              grey
}

type brush func(string) string

type consoleWriter struct {
	lg       *logWriter
	Level    LogLevel `json:"level"`
	Colorful bool     `json:"color"` //this filed is useful only when system's terminal supports color
}

func init() {
	Register(AdapterConsole, NewConsole)
}

func NewConsole() Logger {
	cw := &consoleWriter{
		lg:       newLogWriter(&ansiColorWriter{w: os.Stdout}),
		Level:    LevelDebug,
		Colorful: true,
	}
	return cw
}

func (c *consoleWriter) Init(jsonConfig string) error {
	if len(jsonConfig) == 0 {
		return nil
	}
	err := json.Unmarshal([]byte(jsonConfig), c)
	return err
}

func (c *consoleWriter) Destroy() {

}

func (c *consoleWriter) Flush() {

}

func (c *consoleWriter) WriteMsg(when time.Time, msg string, level LogLevel) error {
	if level > c.Level {
		return nil
	}
	if c.Colorful {
		msg = colors[level](msg)
	}
	c.lg.println(when, msg)
	return nil
}

// newBrush return a fix color Brush
func newBrush(color string) brush {
	pre := "\x1b["
	reset := "\x1b[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}
