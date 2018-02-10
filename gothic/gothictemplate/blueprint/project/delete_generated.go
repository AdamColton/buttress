package project

import (
	"github.com/adamcolton/buttress/config"
	"github.com/adamcolton/buttress/walker"
)

func DeleteGeneratedFiles() {
	w := walker.New()
	w.Dir = walker.DirNo
	err := w.SetMatch(`(.+((\.gen)|(\.gothic))\..+)`)
	if err != nil {
		panic(err)
	}
	w.Callback = func(v *walker.Visit) {
		v.Delete = true
	}
	w.Walk(config.MustGetString("path.app"))
}
