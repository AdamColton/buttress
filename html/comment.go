package html

import (
	"github.com/adamcolton/gothic/gothicio"
	"io"
)

// Comment represents an html text node.
type Comment struct {
	Text string
	Wrap bool
	*parent
}

// NewComment takes the comment string and returns a Comment Node.
func NewComment(text string) *Comment {
	c := &Comment{
		Text: text,
		Wrap: true,
	}
	c.parent = newParent(c)
	return c
}

// Write a Text Node to a writer.
func (c *Comment) WriteTo(w io.Writer) (n int64, err error) {
	nw := newWriter(w)
	c.write(nw)
	return nw.Sum, nw.Err
}

var (
	startComment = []byte("<!--")
	endComment   = []byte("-->")
)

func (c *Comment) write(w writer) {
	w.Write(startComment)
	lww := gothicio.NewLineWrappingWriter(w)
	lww.Write([]byte(c.Text))
	w.Write(endComment)
}
