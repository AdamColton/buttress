package query

import (
	"github.com/adamcolton/buttress/html"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/grammar/regexgram"
	"github.com/adamcolton/parlex/lexer/simplelexer"
	"github.com/adamcolton/parlex/parser/packrat"
	"github.com/adamcolton/parlex/tree"
)

const (
	lexerRules = `
class   /\.[\w\-]+/
id      /#[\w\-]+/
tag     /[\w\-]+/
comma   /[ \t]*,[ \t]*/
child   /[ \t]*>[ \t]*/
nextSib /\+/
after   /\~/
descd   /[ \t]+/
`

	grammarRules = `
Selectors -> (Selector comma)* Selector
Selector  -> Element (child|after|descd|nextSib Element)*
Element   -> tag class* id? class*
          -> class class* id? class*
          -> class* id class*
`
)

var lxr = parlex.MustLexer(simplelexer.New(lexerRules))
var grmr, grmrRdcr = regexgram.Must(grammarRules)
var prsr = packrat.New(grmr)

var rdcr = tree.Merge(grmrRdcr, tree.Reducer{
	"Selectors": tree.RemoveAll("comma"),
})

var runner = parlex.New(lxr, prsr, rdcr)

type nodeChecker interface {
	check(html.TagNode) bool
}

type checkTag string

func (m checkTag) check(node html.TagNode) bool {
	return node.Name() == string(m)
}

type checkClass string

func (m checkClass) check(node html.TagNode) bool {
	for _, class := range html.Classes(node) {
		if class == string(m) {
			return true
		}
	}
	return false
}

type checkID string

func (m checkID) check(node html.TagNode) bool {
	id, _ := node.Attribute("id")
	return id == string(m)
}

type locationChecker func([]int) bool

type descendant []int

func newDescendant(where []int) locationChecker { return descendant(where).check }

// checks that loc is a descendant of where
func (where descendant) check(loc []int) bool {
	if len(where) > len(loc) {
		return false
	}
	for i, w := range where {
		if w != loc[i] {
			return false
		}
	}
	return true
}

type child []int

func newChild(where []int) locationChecker { return child(where).check }

func (where child) check(loc []int) bool {
	if len(where)+1 != len(loc) {
		return false
	}
	for i, w := range where {
		if w != loc[i] {
			return false
		}
	}
	return true
}

type after []int

func newAfter(where []int) locationChecker { return after(where).check }

func (where after) check(loc []int) bool {
	if len(where) != len(loc) {
		return false
	}
	for i, w := range where[:len(where)-1] {
		if w != loc[i] {
			return false
		}
	}
	return where[len(where)-1] < loc[len(loc)-1]
}

type nextSib []int

func newNextSib(where []int) locationChecker { return nextSib(where).check }

func (where nextSib) check(loc []int) bool {
	if len(where) != len(loc) {
		return false
	}
	for i, w := range where[:len(where)-1] {
		if w != loc[i] {
			return false
		}
	}
	return where[len(where)-1]+1 == loc[len(loc)-1]
}

type selector struct {
	loc      locationChecker
	checkers []nodeChecker
	next     *selector
	nextLoc  func([]int) locationChecker
}

func (s *selector) check(node html.TagNode, loc *Location) bool {
	if s.loc != nil && !s.loc(loc.Tag) {
		return false
	}
	for _, c := range s.checkers {
		if !c.check(node) {
			return false
		}
	}
	return true
}

type selectors []*selector

func (s selectors) QueryAll(node html.Node) []html.TagNode {
	op := &selectorsOp{
		selectors: s,
		queryAll:  true,
	}
	Walk(node, op.compare)
	return op.matches
}

func (s selectors) Query(node html.Node) html.TagNode {
	op := &selectorsOp{
		selectors: s,
	}
	Walk(node, op.compare)
	if len(op.matches) == 0 {
		return nil
	}
	return op.matches[0]
}

func (s selectors) QueryAllPaths(node html.Node) Paths {
	op := &selectorsOp{
		selectors: s,
		queryAll:  true,
		path:      true,
	}
	Walk(node, op.compare)
	return op.paths
}

func (s selectors) QueryPath(node html.Node) Path {
	op := &selectorsOp{
		selectors: s,
		path:      true,
	}
	Walk(node, op.compare)
	if len(op.paths) == 0 {
		return nil
	}
	return op.paths[0]
}

type selectorsOp struct {
	selectors
	matches  []html.TagNode
	paths    Paths
	queryAll bool
	path     bool
}

// this is just to organize my thoughts
func (op *selectorsOp) compare(node html.Node, loc *Location) error {
	tag, ok := node.(html.TagNode)
	if !ok {
		return nil
	}

	var toAppend []*selector

	matched := false
	for _, sel := range op.selectors {
		if matched && sel.next == nil {
			// no need to compare final acceptors
			continue
		}
		if sel.check(tag, loc) {
			if sel.next == nil {
				if op.path {
					op.paths = append(op.paths, loc.Path)
				} else {
					op.matches = append(op.matches, tag)
				}
				if !op.queryAll {
					return ErrStop
				}
				matched = true
			} else {
				next := &selector{
					loc:      sel.nextLoc(loc.Tag),
					checkers: sel.next.checkers,
					next:     sel.next.next,
					nextLoc:  sel.nextLoc,
				}
				toAppend = append(toAppend, next)
			}
		}
	}

	op.selectors = append(op.selectors, toAppend...)
	return nil
}

type SelectorQuery interface {
	QueryAll(node html.Node) []html.TagNode
	Query(node html.Node) html.TagNode
	QueryAllPaths(node html.Node) Paths
	QueryPath(node html.Node) Path
}

func Selector(str string) (SelectorQuery, error) {
	parsedSelectors, err := runner.Run(str)
	if err != nil {
		return nil, err
	}

	sels := make(selectors, parsedSelectors.Children())
	for i := range sels {
		sels[i] = parseSelector(parsedSelectors.Child(i))
	}
	return sels, nil
}

func MustSelector(str string) SelectorQuery {
	s, err := Selector(str)
	if err != nil {
		panic(err)
	}
	return s
}

func parseSelector(parsedSelector parlex.ParseNode) *selector {
	var root, cur *selector
	for i := 0; i < parsedSelector.Children(); i++ {
		child := parsedSelector.Child(i)
		switch child.Kind().String() {
		case "Element":
			if root == nil {
				root = parseElement(child)
				cur = root
			} else {
				cur.next = parseElement(child)
				cur = cur.next
			}
		case "descd":
			cur.nextLoc = newDescendant
		case "child":
			cur.nextLoc = newChild
		case "after":
			cur.nextLoc = newAfter
		case "nextSib":
			cur.nextLoc = newNextSib
		}
	}
	return root
}

func parseElement(element parlex.ParseNode) *selector {
	s := &selector{}
	for i := 0; i < element.Children(); i++ {
		checkType := element.Child(i)
		switch checkType.Kind().String() {
		case "tag":
			s.checkers = append(s.checkers, checkTag(checkType.Value()))
		case "class":
			v := []byte(checkType.Value())[1:] // remove leading .
			s.checkers = append(s.checkers, checkClass(v))
		case "id":
			v := []byte(checkType.Value())[1:] // remove leading #
			s.checkers = append(s.checkers, checkID(v))
		}
	}
	return s
}
