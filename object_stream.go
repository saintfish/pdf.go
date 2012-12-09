package pdf

import (
	"bytes"
	"errors"
	"io"
)

var InconsistentStreamLength = errors.New("Length of stream is inconsistent to the length in dictionary")

type writerWrapper Writer

func (w *writerWrapper) Write(b []byte) (n int, err error) {
	return (*Writer)(w).write(b)
}

type Stream struct {
	dict Dictionary
	buf  *bytes.Buffer
}

func NewStream(dict Dictionary) Stream {
	return Stream{
		dict: dict,
		buf:  new(bytes.Buffer),
	}
}

func (s Stream) Write(p []byte) (n int, err error) {
	return s.buf.Write(p)
}

func (s Stream) writeTo(w *Writer) (err error) {
	if length, hasLength := s.dict["Length"]; hasLength && length != s.buf.Len() {
		return InconsistentStreamLength
	}
	s.dict["Length"] = s.buf.Len()
	if err = s.dict.writeTo(w); err != nil {
		return
	}
	if _, err = w.writeString("\nstream\n"); err != nil {
		return
	}
	if _, err = io.Copy((*writerWrapper)(w), s.buf); err != nil {
		return
	}
	_, err = w.writeString("\nendstream")
	return
}
