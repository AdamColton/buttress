package mutate

import (
	"fmt"
	"github.com/adamcolton/buttress/html"
	"github.com/adamcolton/buttress/html/query"
)

type Mutator interface {
	Mutate(html.Node) html.Node
}

type MutateFunc func(html.Node) html.Node

func (fn MutateFunc) Mutate(root html.Node) html.Node {
	return fn(root)
}

type MutateChain []Mutator

func (c *MutateChain) Mutate(node html.Node) html.Node {
	for _, m := range *c {
		node = m.Mutate(node)
	}
	return node
}

func (c *MutateChain) AddMutators(mutators ...Mutator) {
	*c = append(*c, mutators...)
}

// Chain several mutators together. There is special logic that will skip the
// first mutator if it is nil and cast it to a MutateChain if that's the
// underlying type. This makes it easy to append mutations.
func Chain(mutators ...Mutator) Mutator {
	if len(mutators) == 0 {
		return nil
	}
	var c *MutateChain
	if mutators[0] == nil {
		mutators = mutators[1:]
	} else if tc, ok := mutators[0].(*MutateChain); ok {
		c = tc
		mutators = mutators[1:]
	}
	ret := append(*c, mutators...)
	return &ret
}

func AppendTags(selector string, attrs ...string) (MutateFunc, error) {
	q, err := query.Selector(selector)
	if err != nil {
		return nil, err
	}

	return func(root html.Node) html.Node {
		tags := q.QueryAll(root)
		for _, t := range tags {
			for i := 1; i < len(attrs); i += 2 {
				if v, ok := t.Attribute(attrs[i-1]); ok {
					t.AddAttributes(attrs[i-1], fmt.Sprintf("%s %s", v, attrs[i]))
				} else {
					t.AddAttributes(attrs[i-1], attrs[i])
				}
			}
		}
		return root
	}, nil
}

func AppendClass(selector, class string) (MutateFunc, error) {
	q, err := query.Selector(selector)
	if err != nil {
		return nil, err
	}

	return func(root html.Node) html.Node {
		tags := q.QueryAll(root)
		for _, t := range tags {
			t.AppendClass(class)
		}
		return root
	}, nil
}
