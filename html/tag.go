package html

import (
	"io"
	"strings"
)

// Tag is any tag with an opening and closing that may have children
type Tag struct {
	tag string
	attributes
	*fragment
}

// NewTag creates a Tag. It requires a tag name and optionally takes attributes.
func NewTag(tag string, attrs ...string) *Tag {
	t := &Tag{
		tag:        tag,
		attributes: newAttributes(attrs),
		fragment:   &fragment{},
	}
	t.parent = newParent(t)
	return t
}

// Name returns the tag name
func (t *Tag) Name() string { return t.tag }

// Write a tag to a writer
func (t *Tag) WriteTo(w io.Writer) (n int64, err error) {
	nw := newWriter(w)
	t.write(nw)
	return nw.Sum, nw.Err
}

var (
	openBracket  = []byte("<")
	closeBracket = []byte(">")
	openCloseTag = []byte("</")
	closeVoidTag = []byte("/>")
)

func (t *Tag) write(w writer) {
	bTag := []byte(t.tag)
	s0 := w.Sum
	w.Write(openBracket)
	w.Write(bTag)
	closeLen := int(w.Sum-s0) + 2
	t.attributes.write(w)
	w.Write(closeBracket)
	lineLen := int(w.Sum-s0) + len(w.padding)

	multiline := true
	if l := len(t.fragment.children); l == 0 {
		multiline = false
	} else if l == 1 {
		if text, ok := t.fragment.children[0].(*Text); ok {
			multiline = strings.Contains(text.Text, "\n") || len(text.Text)+lineLen+closeLen > w.wrapWidth
		}
	}

	cw := w.indent()
	if multiline {
		cw.nl()
	}
	t.fragment.write(cw)
	if multiline {
		w.nl()
	}
	w.Write(openCloseTag)
	w.Write(bTag)
	w.Write(closeBracket)
}

// VoidTag is a self closing tag that cannot contain children
type VoidTag struct {
	tag string
	attributes
	*parent
}

// NewVoidTag creates a VoidTag. It requires a tag name and optionally takes
// attributes.
func NewVoidTag(tag string, attrs ...string) *VoidTag {
	t := &VoidTag{
		tag:        tag,
		attributes: newAttributes(attrs),
	}
	t.parent = newParent(t)
	return t
}

// Name returns the tag name
func (t *VoidTag) Name() string { return t.tag }

// Write a VoidTag to a writer
func (t *VoidTag) WriteTo(w io.Writer) (n int64, err error) {
	nw := newWriter(w)
	t.write(nw)
	return nw.Sum, nw.Err
}

func (t *VoidTag) write(w writer) {
	w.Write(openBracket)
	w.Write([]byte(t.tag))
	t.attributes.write(w)
	w.Write(closeVoidTag)
}
