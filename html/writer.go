package html

import (
	"github.com/adamcolton/gothic/gothicio"
	"io"
)

var (
	WrapWidth = 80
	Indent    = "  "
)

type writer struct {
	*gothicio.SumWriter
	wrapWidth int
	padding   string
	indentStr string
}

func newWriter(w io.Writer) writer {
	if w, ok := w.(writer); ok {
		return w
	}
	sw, ok := w.(*gothicio.SumWriter)
	if !ok {
		sw = &gothicio.SumWriter{
			Writer: w,
		}
	}
	return writer{
		SumWriter: sw,
		wrapWidth: WrapWidth,
		indentStr: Indent,
	}
}

func (w writer) WrapWidth() int  { return w.wrapWidth }
func (w writer) Padding() string { return w.padding }

func (w writer) indent() writer {
	w.padding += w.indentStr
	return w
}

var nl = []byte("\n")

func (w writer) nl() {
	w.Write(nl)
	w.Write([]byte(w.padding))
}
