package main

import (
	"github.com/adamcolton/gothic/gothicgo"
)

func setTypes() {
	typeSet{
		{
			name:   "primary",
			gotype: gothicgo.Uint64Type,
		},
	}.set()
}

type typeDef struct {
	name   string
	gotype gothicgo.Type
}

func (td typeDef) set() {
	if td.gotype != nil {
		gomodel.Types[ts.name] = ts.gotype
	}
}

type typeSet []typeDef

func (t typeSet) Set() {
	for _, td := range t {
		td.Set()
	}
}
