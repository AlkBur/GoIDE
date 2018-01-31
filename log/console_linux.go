package log

type ansiColorWriter struct {
	w    io.Writer
}

func (cw *ansiColorWriter) Write(p []byte) (int, error) {
	return cw.w.Write(p)
}

// newBrush return a fix color Brush
func newBrush(color string) brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}
