package pdf

type Dictionary map[Name]Object

func (d Dictionary) writeTo(w *Writer) (err error) {
	if _, err = w.writeString("<<\n"); err != nil {
		return
	}
	for name, value := range d {
		if err = name.writeTo(w); err != nil {
			return
		}
		if _, err = w.write([]byte{' '}); err != nil {
			return
		}
		if err = writeObject(value, w); err != nil {
			return
		}
		if _, err = w.write([]byte{'\n'}); err != nil {
			return
		}
	}
	_, err = w.writeString(">>")
	return
}
