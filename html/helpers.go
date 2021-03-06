package html

import (
	"bytes"
	"strings"
)

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
