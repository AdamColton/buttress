package html

import (
	"io"
)

// Fragment is a container for holding html Nodes. Fragments can be useful when
// constructing a document but bee careful when using it, it behaves like a Node
// but doesn't render anything so the children of a Fragment will actually be
// renders as "siblings" to their "uncles". To avoid this, use RemoveFragments
// before traversing the tree.
type Fragment struct {
	*fragment
}

// NewFragment makes a new Fragment. It takes child Nodes as a variadic arg.
func NewFragment(children ...Node) *Fragment {
	f := &Fragment{
		fragment: &fragment{
			children: children,
		},
	}
	f.parent = newParent(f)
	return f
}

type fragment struct {
	children []Node
	*parent
}

func (f *fragment) WriteTo(w io.Writer) (n int64, err error) {
	nw := newWriter(w)
	f.write(nw)
	return int64(nw.Sum), nw.Err
}

func (f *fragment) frag() *fragment {
	return f
}

func (f *fragment) write(w writer) {
	pad := false
	for _, c := range f.children {
		if pad {
			w.nl()
		} else {
			pad = true
		}
		c.write(w)
	}
}

func (f *fragment) Children() int { return len(f.children) }

func (f *fragment) Child(idx int) Node {
	if idx < 0 {
		idx += len(f.children)
	}
	if idx < 0 || idx >= len(f.children) {
		return nil
	}
	return f.children[idx]
}

func (f *fragment) AddChildren(children ...Node) {
	// TODO: don't allow add child without nil parent
	for _, c := range children {
		c.setParent(f.self)
	}
	f.children = append(f.children, children...)
}

func (f *fragment) RemoveChild(idx int) {
	if idx < 0 {
		idx += len(f.children)
	}
	if idx < 0 || idx >= len(f.children) {
		return
	}
	f.children[idx].setParent(nil)
	if idx == len(f.children)-1 {
		f.children = f.children[0:idx]
	} else {
		f.children = append(f.children[0:idx], f.children[idx+1:]...)
	}
}

func (f *fragment) RemoveAll() {
	for _, c := range f.children {
		c.setParent(nil)
	}
	f.children = nil
}

func (f *fragment) ReplaceChildren(frag *Fragment) {
	f.replaceChildren(frag.fragment)
}

func (f *fragment) replaceChildren(frag *fragment) {
	f.children = frag.children
	for _, c := range f.children {
		c.setParent(f)
	}
	frag.children = nil
}

func (f *fragment) Text(text string) *Text {
	t := NewText(text)
	f.AddChildren(t)
	return t
}

func (f *fragment) Tag(tag string, attrs ...string) *Tag {
	t := NewTag(tag, attrs...)
	f.AddChildren(t)
	return t
}

func (f *fragment) VoidTag(tag string, attrs ...string) *VoidTag {
	t := NewVoidTag(tag, attrs...)
	f.AddChildren(t)
	return t
}

// RemoveFragments modifies a tree and removes any Fragments other than the
// root.
func RemoveFragments(node Node) {
	removeFragments(node, 0)
}

func removeFragments(node Node, idx int) {
	if p := node.Parent(); p != nil {
		if f, ok := node.(*Fragment); ok {
			for _, c := range f.children {
				c.setParent(p)
			}
			pf := p.(ContainerNode).frag()
			if idx == len(pf.children)-1 {
				pf.children = append(pf.children[:idx], f.children...)
			} else {
				pf.children = append(pf.children[:idx], append(f.children, pf.children[idx+1:]...)...)
			}
			f.parent = nil
			f.children = nil
		}
	}

	if c, ok := node.(ContainerNode); ok {
		for i := 0; i < c.Children(); i++ {
			removeFragments(c.Child(i), i)
		}
	}
}
