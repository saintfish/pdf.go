package pdf

import (
	"errors"
	"io"
	"unicode"
)

func escapeAndWriteByte(b byte, w io.Writer) (err error) {
	switch b {
	case '\r':
		_, err = w.Write([]byte("\\r"))
	case '\n':
		_, err = w.Write([]byte("\\n"))
	case '\t':
		_, err = w.Write([]byte("\\t"))
	case '\b':
		_, err = w.Write([]byte("\\b"))
	case '\f':
		_, err = w.Write([]byte("\\f"))
	case '(', ')', '\\':
		_, err = w.Write([]byte{'\\', b})
	default:
		_, err = w.Write([]byte{b})
	}
	return
}

func escapeAndWriteBytes(bs []byte, w io.Writer) (err error) {
	for _, b := range bs {
		if err = escapeAndWriteByte(b, w); err != nil {
			return
		}
	}
	return
}

var EncodeError = errors.New("Unable to encode string with encoder")

type StringEncoding interface {
	Encode(str string, w io.Writer) error
}

type asciiEncoding struct {
}

func (*asciiEncoding) Encode(str string, w io.Writer) (err error) {
	if _, err = w.Write([]byte{'('}); err != nil {
		return
	}
	for _, r := range str {
		if r > unicode.MaxASCII {
			return EncodeError
		}
		if err = escapeAndWriteByte(byte(r&0xFF), w); err != nil {
			return
		}
	}
	_, err = w.Write([]byte{')'})
	return
}

type utf16BeEncoding struct {
}

func (*utf16BeEncoding) Encode(str string, w io.Writer) (err error) {
	w.Write([]byte{'('})
	if err = escapeAndWriteBytes([]byte{0xFE, 0xFF}, w); err != nil {
		return
	}
	b := make([]byte, 2)
	for _, r := range str {
		b[0] = byte(r >> 8)
		b[1] = byte(r & 0xFF)
		if err = escapeAndWriteBytes(b, w); err != nil {
			return
		}
	}
	w.Write([]byte{')'})
	return nil
}

type hexEncoding struct {
}

const hexString = "0123456789ABCDEF"

func (*hexEncoding) Encode(str string, w io.Writer) (err error) {
	if _, err = w.Write([]byte{'<'}); err != nil {
		return
	}
	b := make([]byte, 2)
	for i := 0; i < len(str); i++ {
		var c byte = str[i]
		b[0] = hexString[c>>4]
		b[1] = hexString[c&0xFF]
		if _, err = w.Write(b); err != nil {
			return
		}
	}
	_, err = w.Write([]byte{'>'})
	return
}

var (
	AsciiEncoding   StringEncoding = new(asciiEncoding)
	Utf16BeEncoding StringEncoding = new(utf16BeEncoding)
	// PdfDocEncoding
	HexEncoding StringEncoding = new(hexEncoding)
)

type String struct {
	str      string
	encoding StringEncoding
}

func NewString(str string) String {
	return String{
		str:      str,
		encoding: AsciiEncoding,
	}
}

func (s String) writeTo(w *Writer) error {
	return s.encoding.Encode(s.str, w.w)
}
