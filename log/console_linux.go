package log

type ansiColorWriter struct {
	w io.Writer
}

func (cw *ansiColorWriter) Write(p []byte) (int, error) {
	return cw.w.Write(p)
}
