package html

import (
	"io"
)

// Doctype represents a doctype tag. It does not fulfil TagNode.
type Doctype struct {
	doctype string
	*parent
}

// NewDoctype returns a new doctype tag. To create a standard doctype tag, just
// call NewDoctype("html")
func NewDoctype(doctype string) *Doctype {
	d := &Doctype{
		doctype: doctype,
	}
	d.parent = newParent(d)
	return d
}

// WriteTo Doctype to an io.Writer
func (d *Doctype) WriteTo(w io.Writer) (n int64, err error) {
	nw := newWriter(w)
	d.write(nw)
	return nw.Sum, nw.Err
}

var (
	openDocType = []byte("<!DOCTYPE ")
)

func (d *Doctype) write(w writer) {
	w.Write(openDocType)
	w.Write([]byte(d.doctype))
	w.Write(closeBracket)
}
