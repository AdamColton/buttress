package bootstrap3

import (
	"fmt"
	"github.com/adamcolton/buttress/bootstrap3/csssize"
	"github.com/adamcolton/buttress/html"
	"strings"
)

// https://getbootstrap.com/docs/3.3/css/#grid

type Container struct {
	Fluid bool
	rows  []*Row
}

func NewContainer() *Container {
	return &Container{}
}

func NewFluidContainer() *Container {
	return &Container{
		Fluid: true,
	}
}

func (c *Container) NewRow(cells ...*Cell) *Row {
	r := NewRow(html.NewTag("div", "class", "row"))
	c.rows = append(c.rows, r)
	r.AddCells(cells...)
	return r
}

func (c *Container) AddNodes(class ColClass, nodes ...html.Node) {
	for _, node := range nodes {
		cell := &Cell{
			Class:    class,
			Contents: node,
		}
		if len(c.rows) == 0 {
			c.NewRow(cell)
			continue
		}
		c.rows[len(c.rows)-1].AddCells(cell)
	}
}

func (c *Container) AddCells(cells ...*Cell) {
	if len(c.rows) == 0 {
		c.NewRow(cells...)
		return
	}
	c.rows[len(c.rows)-1].AddCells(cells...)
}

func (c *Container) Render() html.Node {
	class := "container"
	if c.Fluid {
		class += "-fluid"
	}
	div := html.NewTag("div", "class", class)
	for _, r := range c.rows {
		div.AddChildren(r.Render())
	}
	return div
}

type Row struct {
	container html.ContainerNode
	cells     []*Cell
}

func NewRow(container html.ContainerNode) *Row {
	return &Row{
		container: container,
	}
}

func (r *Row) AddNode(class ColClass, node html.Node) {
	r.cells = append(r.cells, &Cell{
		Class:    class,
		Contents: node,
	})
}

func (r *Row) AddCells(cells ...*Cell) {
	r.cells = append(r.cells, cells...)
}

func (r *Row) Render() html.Node {
	if r.container == nil {
		r.container = html.NewFragment()
	}
	counts := make(map[string]int)
	for _, cell := range r.cells {
		var clear []string
		for _, size := range csssize.Sizes() {
			col, _ := cell.Class.Size(size)
			offset, _ := cell.Class.Offset(size)
			total := col + offset
			if cell.Class.Hidden(size) {
				total = 0
			}
			if total+counts[size.String()] > 12 {
				clear = append(clear, size.String())
				counts[size.String()] = total
			} else {
				counts[size.String()] += total
			}
		}
		if len(clear) != 0 {
			class := make([]string, len(clear)+1)
			class[0] = "clearfix"
			for i, size := range clear {
				class[i+1] = fmt.Sprintf("visible-%s-block", size)
			}
			r.container.AddChildren(html.NewTag("div", "class", strings.Join(class, " ")))
		}
		r.container.AddChildren(cell.Render())
	}
	return r.container
}

type Cell struct {
	Class    ColClass
	Contents html.Node
}

func (c *Cell) Render() html.Node {
	div := html.NewTag("div", "class", c.Class.String())
	div.AddChildren(c.Contents)
	return div
}

// TODO: https://getbootstrap.com/docs/3.3/css/#grid-responsive-resets
// Keep track of the columns in a row for each size and insert clearfix where
// needed.
