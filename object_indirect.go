package pdf

import (
	"fmt"
)

type Indirect struct {
	Id  int32
	Gen int32
}

func (i Indirect) writeTo(w *Writer) (err error) {
	_, err = w.writeString(fmt.Sprintf("%d %d R", i.Id, i.Gen))
	return
}
