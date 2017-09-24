package bootstrap3

import (
	"fmt"
	"github.com/adamcolton/buttress/bootstrap3/csssize"
	"strings"
)

type ColClass interface {
	SetSize(size int, class csssize.CSSSize) error
	Size(class csssize.CSSSize) (int, bool)
	SetOffset(offset int, class csssize.CSSSize) error
	Offset(class csssize.CSSSize) (int, bool)
	SetHidden(isHidden bool, class csssize.CSSSize)
	Hidden(class csssize.CSSSize) bool
	SetSizeOffset(size, offset int, class csssize.CSSSize) error
	String() string
}

func NewColClass() ColClass {
	return &colClass{
		offsets: make(map[string]int),
		sizes:   make(map[string]int),
		hidden:  make(map[string]bool),
	}
}

type colClass struct {
	offsets map[string]int
	sizes   map[string]int
	hidden  map[string]bool
	str     string
}

func (c *colClass) String() string {
	if c.str != "" {
		return c.str
	}
	var class []string
	var hide []string
	var visible string
	lastCol, lastOffset := 0, 0
	for _, size := range csssize.Sizes() {
		sizeStr := size.String()
		col, hidden := c.sizes[sizeStr], c.hidden[sizeStr]
		offset, offsetSet := c.offsets[sizeStr]
		if col != 0 && col != lastCol {
			class = append(class, fmt.Sprintf("col-%s-%d", sizeStr, col))
			lastCol = col
		}
		if offset != lastOffset && (offsetSet || offset != 0) {
			class = append(class, fmt.Sprintf("col-%s-offset-%d", sizeStr, offset))
			lastOffset = offset
		}
		if hidden {
			hide = append(hide, fmt.Sprintf("hidden-%s", sizeStr))
		} else {
			visible = sizeStr
		}
	}
	if len(hide) == 3 {
		class = append(class, fmt.Sprintf("visible-%s", visible))
	} else {
		class = append(class, hide...)
	}
	c.str = strings.Join(class, " ")
	return c.str
}

func set(field, key string, val int, m map[string]int) error {
	if val > 12 {
		return fmt.Errorf("%d cannot be larger than 12, got %d", field, val)
	}
	if val < 0 {
		return fmt.Errorf("%d cannot be less than 0, got %d", field, val)
	}
	m[key] = val
	return nil
}

func get(class csssize.CSSSize, m map[string]int) (int, bool) {
	if val, defined := m[class.String()]; defined {
		return val, true
	}
	for class = class.Smaller(); class != nil; class = class.Smaller() {
		if val, defined := m[class.String()]; defined {
			return val, false
		}
	}
	return 0, false
}

func (c *colClass) SetSize(size int, class csssize.CSSSize) error {
	c.str = ""
	return set("Size", class.String(), size, c.sizes)
}

func (c *colClass) Size(class csssize.CSSSize) (int, bool) {
	return get(class, c.sizes)
}

func (c *colClass) SetOffset(offset int, class csssize.CSSSize) error {
	c.str = ""
	return set("Offset", class.String(), offset, c.offsets)
}

func (c *colClass) Offset(class csssize.CSSSize) (int, bool) {
	return get(class, c.offsets)
}

func (c *colClass) SetHidden(isHidden bool, class csssize.CSSSize) {
	c.str = ""
	c.hidden[class.String()] = isHidden
}

func (c *colClass) Hidden(class csssize.CSSSize) bool {
	return c.hidden[class.String()]
}

func (c *colClass) SetSizeOffset(size, offset int, class csssize.CSSSize) error {
	c.str = ""
	err := set("Size", class.String(), size, c.sizes)
	if err != nil {
		return err
	}
	err = set("Offset", class.String(), offset, c.offsets)
	return err
}
