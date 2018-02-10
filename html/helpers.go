package html

import (
	"bytes"
	"io"
	"strings"
)

// NewLine is the string that will be used when rendering HTML to a writer for a
// newline.
var NewLine = "\n"

// Padding is the string that is used to indent html
var Padding = "  "

type stringWriter interface {
	WriteString(s string) (n int, err error)
}

// StringWriterWrapper can wrap any io.Writer to fulfill StringWriter
type stringWriterWrapper struct {
	io.Writer
}

// WriteString fulfils StringWriter on StringWriterWrapper. It just casts the
// string to []byte
func (w stringWriterWrapper) WriteString(s string) (n int, err error) {
	return w.Write([]byte(s))
}

func toStringWriter(w io.Writer) stringWriter {
	if sw, ok := w.(stringWriter); ok {
		return sw
	}
	return stringWriterWrapper{w}
}

type writer struct {
	sw            stringWriter
	onNewLine     bool
	start         int
	padding       string
	parentPadding string
	*counter
}

type counter struct {
	err error
	sum int
}

func newWriter(w io.Writer) *writer {
	return &writer{
		sw:      toStringWriter(w),
		counter: &counter{},
	}
}

func (w *writer) write(str string) {
	if w.err != nil {
		return
	}
	w.onNewLine = false
	n, err := w.sw.WriteString(str)
	w.sum += n
	w.err = err
}

func (w *writer) nl() {
	if w.err != nil {
		return
	}
	w.onNewLine = true
	n, err := w.sw.WriteString(NewLine)
	w.sum += n
	w.err = err
	if w.err != nil {
		return
	}
	n, err = w.sw.WriteString(w.padding)
	w.sum += n
	w.err = err
}

func (w *writer) pnl() {
	w.onNewLine = true
	n, err := w.sw.WriteString(NewLine)
	w.sum += n
	w.err = err
	if w.err != nil {
		return
	}
	n, err = w.sw.WriteString(w.parentPadding)
	w.sum += n
	w.err = err
}

func (w *writer) inc() *writer {
	cp := *w
	cp.parentPadding = cp.padding
	cp.padding += Padding
	return &cp
}

// String uses bytes.Buffer to render html as a string. This tends to be useful
// in testing, but less so in production code.
func String(node Node) string {
	var buf bytes.Buffer
	node.write(newWriter(&buf))
	return buf.String()
}

// Classes is a helper function that returns the classes on a TagNode as a slice
// of strings.
func Classes(node TagNode) []string {
	classes, _ := node.Attribute("class")
	return strings.Fields(classes)
}

// parent is a helper that is embeded to hold a reference to self and parent and
// fulfill the Parent() and setParent() methods on Node.
type parent struct {
	parent Node
	self   Node
}

func newParent(self Node) *parent {
	return &parent{self: self}
}

func (p *parent) Parent() Node {
	return p.parent
}

func (p *parent) setParent(newparent Node) {
	p.parent = newparent
}
