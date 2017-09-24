package query

import (
	"github.com/adamcolton/buttress/html"
)

type Prototype struct {
	html      html.Node
	queries   map[string]Path
	queryAlls map[string]Paths
}

func NewPrototype(html html.Node) *Prototype {
	return &Prototype{
		html:      html,
		queries:   make(map[string]Path),
		queryAlls: make(map[string]Paths),
	}
}

func (p *Prototype) Queries(queries ...string) *Prototype {
	for i := 1; i < len(queries); i += 2 {
		p.queries[queries[i-1]] = MustSelector(queries[i]).QueryPath(p.html)
	}
	return p
}

func (p *Prototype) QueryAlls(queries ...string) *Prototype {
	for i := 1; i < len(queries); i += 2 {
		p.queryAlls[queries[i-1]] = MustSelector(queries[i]).QueryAllPaths(p.html)
	}
	return p
}

func (p *Prototype) Clone() html.Node {
	return html.Clone(p.html)
}

func (p *Prototype) Query(name string, node html.Node) html.Node {
	return p.queries[name].Query(node)
}

func (p *Prototype) QueryAll(name string, node html.Node) []html.Node {
	return p.queryAlls[name].Query(node)
}
