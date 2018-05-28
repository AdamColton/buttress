package mutate

import (
	"fmt"
	"github.com/adamcolton/buttress/html"
	"github.com/adamcolton/buttress/html/query"
)

//TODO: much of this stutters:
// MutateChain
// MutateFunc

// Mutator takes in a node, mutates it and returns the mutated version.
type Mutator interface {
	Mutate(html.Node) html.Node
}

// MutateFunc wraps a func that performs a mutation so it fulfils the Mutator
// interface.
type MutateFunc func(html.Node) html.Node

// Mutate fulfils the Mutator interface.
func (fn MutateFunc) Mutate(root html.Node) html.Node {
	return fn(root)
}

// MutateChain applies a chain of mutations to an html Node. It can be used to
// make the same set of changes to many nodes or be embeded into a Node
// generator allowing a set of mutations to be applied after the Node is
// generated.
type MutateChain struct {
	Mutators []Mutator
}

// Mutate applies the MutateChain to the node, also fulfills the Mutator
// interface.
func (c MutateChain) Mutate(node html.Node) html.Node {
	for _, m := range c.Mutators {
		node = m.Mutate(node)
	}
	return node
}

// AddMutators to the chain.
func (c *MutateChain) AddMutators(mutators ...Mutator) {
	c.Mutators = append(c.Mutators, mutators...)
}

// AddMutator allows the output of a Mutator generator (like AppendAttrs or
// AppendClass)
func (c *MutateChain) AddMutator(mutator Mutator, err error) error {
	if err != nil {
		return err
	}
	c.Mutators = append(c.Mutators, mutator)
	return nil
}

// Chain several mutators together. There is special logic that will skip the
// first mutator if it is nil. This makes it easy to append mutations.
func Chain(mutators ...Mutator) Mutator {
	if len(mutators) == 0 {
		return nil
	}

	if mutators[0] == nil {
		mutators = mutators[1:]
		if len(mutators) == 0 {
			return nil
		}
	}

	return MutateChain{
		Mutators: mutators,
	}
}

// AppendAttrs appends the class to the nodes that match the selector. If the
// selector is empty, the class will be appended to the root node.
func AppendAttrs(selector string, attrs ...string) (MutateFunc, error) {
	if selector == "" {
		return func(root html.Node) html.Node {
			if tag, ok := root.(html.TagNode); ok {
				for i := 1; i < len(attrs); i += 2 {
					if v, ok := tag.Attribute(attrs[i-1]); ok {
						tag.AddAttributes(attrs[i-1], fmt.Sprintf("%s %s", v, attrs[i]))
					} else {
						tag.AddAttributes(attrs[i-1], attrs[i])
					}
				}
			}
			return root
		}, nil
	}

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

// AppendClass appends the class to the nodes that match the selector. If the
// selector is empty, the class will be appended to the root node.
func AppendClass(selector, class string) (MutateFunc, error) {
	if selector == "" {
		return func(root html.Node) html.Node {
			if tag, ok := root.(html.TagNode); ok {
				tag.AppendClass(class)
			}
			return root
		}, nil
	}

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
