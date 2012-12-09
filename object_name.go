package pdf

type Name string

func (n Name) writeTo(w *Writer) (err error) {
	if _, err = w.write([]byte{'/'}); err != nil {
		return
	}
	b := make([]byte, 3)
	b[0] = '#'
	for i := 0; i < len(n); i++ {
		c := n[i]
		switch c {
		case ' ', '%', '(', ')', '[', ']', '{', '}', '/', '#':
			b[1] = hexString[c>>4]
			b[2] = hexString[c&0xFF]
			if _, err = w.write(b); err != nil {
				return
			}
		default:
			if c >= 32 && c <= 126 {
				if _, err = w.write([]byte{c}); err != nil {
					return
				}
			} else {
				b[1] = hexString[c>>4]
				b[2] = hexString[c&0xFF]
				if _, err = w.write(b); err != nil {
					return
				}
			}
		}
	}
	return nil
}
