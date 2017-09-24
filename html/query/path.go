package query

import (
	"github.com/adamcolton/buttress/html"
)

type Path []int

func NewPath(ps ...int) Path {
	return Path(ps)
}

func (p Path) Query(n html.Node) html.Node {
	cur := n

	for _, idx := range p {
		parent, ok := cur.(html.ContainerNode)
		if !ok {
			return nil
		}
		cur = parent.Child(idx)
	}

	return cur
}

func (p Path) QueryTag(n html.Node) *html.Tag {
	n = p.Query(n)
	if tag, ok := n.(*html.Tag); ok {
		return tag
	}
	return nil
}

func (p Path) QueryVoidTag(n html.Node) *html.VoidTag {
	n = p.Query(n)
	if tag, ok := n.(*html.VoidTag); ok {
		return tag
	}
	return nil
}

func (p Path) Clone() Path {
	cln := make(Path, len(p))
	copy(cln, p)
	return cln
}

type Paths []Path

func (ps Paths) Query(n html.Node) []html.Node {
	nodes := make([]html.Node, 0, len(ps))
	for _, p := range ps {
		if n := p.Query(n); n != nil {
			nodes = append(nodes, n)
		}
	}
	return nodes
}
