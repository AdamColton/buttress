package html

import (
	"github.com/adamcolton/gothic/gothicio"
	"io"
)

// Text represents an html text node.
type Text struct {
	Text string
	Wrap bool
	*parent
}

// NewText takes the text string and returns a Text Node.
func NewText(text string) *Text {
	t := &Text{
		Text: text,
		Wrap: true,
	}
	t.parent = newParent(t)
	return t
}

// Write a Text Node to a writer.
func (t *Text) WriteTo(w io.Writer) (n int64, err error) {
	nw := newWriter(w)
	t.write(nw)
	return nw.Sum, nw.Err
}

func (t *Text) write(w writer) {
	lww := gothicio.NewLineWrappingWriter(w)
	lww.Write([]byte(t.Text))
}
