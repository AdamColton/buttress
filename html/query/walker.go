package query

import (
	"github.com/adamcolton/buttress/html"
)

// Location of a node with in a document. Path is the path within the data
// structure, which includes Fragments. Node is the node position within
// document, which does not include Fragments. And Tag is the tag position which
// only includes Tags and VoidTags.
type Location struct {
	Node []int
	Tag  []int
	Path Path
}

func (l *Location) Copy() *Location {
	cp := &Location{}
	if l.Node != nil {
		cp.Node = make([]int, len(l.Node))
		copy(cp.Node, l.Node)
	}
	if l.Tag != nil {
		cp.Tag = make([]int, len(l.Tag))
		copy(cp.Tag, l.Tag)
	}
	if l.Path != nil {
		cp.Path = make(Path, len(l.Path))
		copy(cp.Path, l.Path)
	}
	return cp
}

type walkOp struct {
	loc   *Location
	cur   html.ContainerNode
	stack []html.ContainerNode
}

func (op *walkOp) appendPath() { op.loc.Path = append(op.loc.Path, 0) }
func (op *walkOp) appendNode() { op.loc.Node = append(op.loc.Node, 0) }
func (op *walkOp) appendTag()  { op.loc.Tag = append(op.loc.Tag, 0) }
func (op *walkOp) popPath()    { op.loc.Path = op.loc.Path[:len(op.loc.Path)-1] }
func (op *walkOp) popNode()    { op.loc.Node = op.loc.Node[:len(op.loc.Node)-1] }
func (op *walkOp) popTag()     { op.loc.Tag = op.loc.Tag[:len(op.loc.Tag)-1] }
func (op *walkOp) incPath()    { op.loc.Path[len(op.loc.Path)-1]++ }
func (op *walkOp) incNode()    { op.loc.Node[len(op.loc.Node)-1]++ }
func (op *walkOp) incTag()     { op.loc.Tag[len(op.loc.Tag)-1]++ }

// pop tries to pop the stack and returns a bool indicating if it was
// successful. Failure to pop means traversal is done. Popping the stack also
// means that the path needs to be popped and incremented as does the Doc if we
// were traversing a Tag.
func (op *walkOp) pop() bool {
	if len(op.stack) == 0 {
		return false
	}
	if _, ok := op.cur.(*html.Tag); ok {
		op.popTag()
		op.incTag()
		op.popNode()
		op.incNode()
	}
	op.popPath()
	op.incPath()
	op.cur, op.stack = op.stack[len(op.stack)-1], op.stack[:len(op.stack)-1]
	return true
}

func (op *walkOp) push(node html.ContainerNode) {
	op.stack = append(op.stack, op.cur)
	op.cur = node
	op.appendPath()
}

type errStop struct{}

func (errStop) Error() string { return "STOP" }

var ErrStop = errStop{}

// Walk will visit every node and call the corresponding func.
func Walk(root html.Node, handler func(node html.Node, location *Location) error) error {
	op := &walkOp{
		loc: &Location{},
	}
	err := handler(root, op.loc.Copy())
	if err != nil {
		if err == ErrStop {
			return nil
		}
		return err
	}

	// if the root isn't a container, there's nothing else to do
	var ok bool
	op.cur, ok = root.(html.ContainerNode)
	if !ok {
		return nil
	}

	// for the root, we index the children no matter the type
	op.appendPath()
	op.appendTag()
	op.appendNode()

	for {
		idx := op.loc.Path[len(op.loc.Path)-1]
		if idx >= op.cur.Children() {
			if !op.pop() {
				break
			}
			continue
		}

		child := op.cur.Child(idx)
		err := handler(child, op.loc.Copy())
		if err != nil {
			if err == ErrStop {
				return nil
			}
			return err
		}

		switch t := child.(type) {
		case *html.Fragment:
			op.push(t)
		case *html.Tag:
			op.push(t)
			op.appendNode()
			op.appendTag()
		case *html.VoidTag:
			op.incTag()
			op.incNode()
			op.incPath()
		default:
			op.incNode()
			op.incPath()
		}
	}
	return nil
}
