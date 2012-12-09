package pdf

import (
	"fmt"
	"strconv"
)

type Object interface {
}

type pdfWritable interface {
	writeTo(w *Writer) error
}

var (
	nullLiteral = []byte("null")
)

func writeObject(obj Object, w *Writer) (err error) {
	switch o := obj.(type) {
	case nil:
		_, err = w.write(nullLiteral)
	case string:
		err = NewString(o).writeTo(w)
	case int:
		_, err = w.writeString(strconv.Itoa(o))
	case int64:
		_, err = w.writeString(strconv.FormatInt(o, 10))
	case float32:
		_, err = w.writeString(strconv.FormatFloat(float64(o), 'f', -1, 32))
	case float64:
		_, err = w.writeString(strconv.FormatFloat(o, 'f', -1, 64))
	case pdfWritable:
		err = o.writeTo(w)
	default:
		panic(fmt.Sprintf("Unsupported value: %#v <%T>", obj, obj))
	}
	return
}
