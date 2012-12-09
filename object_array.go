package pdf

type Array []Object

func (a Array) writeTo(w *Writer) (err error) {
	if _, err = w.write([]byte{'[', ' '}); err != nil {
		return
	}
	for _, o := range a {
		if err = writeObject(o, w); err != nil {
			return
		}
		if _, err = w.write([]byte{' '}); err != nil {
			return
		}
	}
	_, err = w.write([]byte{']'})
	return
}
