package pdf

import (
	"io"
)

type Writer struct {
	offset int
	w      io.Writer
}

func (w *Writer) write(b []byte) (n int, err error) {
	n, err = w.w.Write(b)
	w.offset += n
	return
}

func (w *Writer) writeString(s string) (n int, err error) {
	return w.write([]byte(s))
}
