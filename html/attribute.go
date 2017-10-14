package html

import (
	"sort"
	"strings"
)

type attributes map[string]string

func newAttributes(attrs []string) attributes {
	a := make(attributes)
	a.AddAttributes(attrs...)
	return a
}

func (a attributes) AddAttributes(attrs ...string) {
	for i := 1; i < len(attrs); i += 2 {
		key := strings.ToLower(attrs[i-1])
		attr := attrs[i]
		if attr != "" {
			a[key] = attr
		}
	}
}

func (a attributes) AddEmptyAttributes(attrs ...string) {
	for _, attr := range attrs {
		a[attr] = ""
	}
}

func (a attributes) Attributes() []string {
	keys := make([]string, len(a))
	i := 0
	for k := range a {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func (a attributes) Attribute(key string) (string, bool) {
	v, ok := a[key]
	return v, ok
}

func (a attributes) write(w *writer) {
	for _, k := range a.Attributes() {
		v := a[k]
		w.write(" ")
		w.write(k)
		if v != "" {
			w.write(`="`)
			w.write(v)
			w.write(`"`)
		}
	}
}

func (a attributes) Remove(key string) {
	key = strings.ToLower(key)
	delete(a, key)
}

func (a attributes) AppendClass(class string) {
	old := a["class"]
	if old == "" {
		a["class"] = class
	} else {
		a["class"] = old + " " + class
	}
}
