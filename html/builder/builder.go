package builder

import (
	"github.com/adamcolton/buttress/html"
)

type Builder struct {
	root  *html.Fragment
	cur   html.ContainerNode
	stack []html.ContainerNode
}

func New() *Builder {
	f := html.NewFragment()
	return &Builder{
		cur:  f,
		root: f,
	}
}

// TODO: func Build() - takes a container node as the root

func (b *Builder) Text(text string) *Builder {
	if text != "" {
		b.cur.AddChildren(html.NewText(text))
	}
	return b
}

func (b *Builder) Tag(tag string, attrs ...string) *Builder {
	t := html.NewTag(tag, attrs...)
	b.cur.AddChildren(t)
	b.push(t)
	return b
}

func (b *Builder) VoidTag(tag string, attrs ...string) *Builder {
	b.cur.AddChildren(html.NewVoidTag(tag, attrs...))
	return b
}

// Close the current container
func (b *Builder) Close() *Builder {
	if l := len(b.stack); l > 0 {
		b.cur = b.stack[l-1]
		b.stack = b.stack[:l-1]
	}
	return b
}

// CloseTag closes the current container and adds additional checks that it is
// a tag and matches the tagStr. CloseTag will panic if these conditions are not
// met.
func (b *Builder) CloseTag(tagStr string) *Builder {
	l := len(b.stack)
	if l == 0 {
		panic("Cannot close, no open tags")
	}
	tag, ok := b.cur.(*html.Tag)
	if !ok {
		panic("Open container is not a tag")
	}
	if tag.Name() != tagStr {
		panic("Tag check failure, expected: " + tagStr + " got:" + tag.Name())
	}
	b.cur = b.stack[l-1]
	b.stack = b.stack[:l-1]
	return b
}

func (b *Builder) push(node html.ContainerNode) {
	b.stack = append(b.stack, b.cur)
	b.cur = node
}

func (b *Builder) Cur() html.ContainerNode {
	return b.cur
}

func (b *Builder) Root() html.Node {
	if b.root.Children() == 1 {
		return b.root.Child(0)
	}
	return b.root
}

func (b *Builder) AddChildren(nodes ...html.Node) {
	b.cur.AddChildren(nodes...)
}
