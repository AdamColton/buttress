package bootstrap3model

import (
	"github.com/adamcolton/buttress/bootstrap3"
	"github.com/adamcolton/gothic/gothicmodel"
)

type FormBuilder struct {
	ValueBuilder ValueBuilder
	Style        *bootstrap3.FormStyle
	ShowPrimary  bool
}

func NewFormBuilder(formStyle *bootstrap3.FormStyle) *FormBuilder {
	return &FormBuilder{
		Style: formStyle,
	}
}

func (fg *FormBuilder) NewForm(fields []gothicmodel.Field, types TypeMap) *bootstrap3.Form {
	var ig InputBuilder
	f := bootstrap3.NewForm(fg.Style)

	for _, field := range fields {
		if !fg.ShowPrimary && field.Primary() {
			ig = HiddenTypeInput
		} else {
			var ok bool
			if types != nil {
				ig, ok = types[field.Type()]
			}
			if !ok {
				ig, ok = Types[field.Type()]
				if !ok {
					continue
				}
			}
		}
		f.AddInputs(ig(field.Name(), fg.ValueBuilder))
	}

	return f
}
