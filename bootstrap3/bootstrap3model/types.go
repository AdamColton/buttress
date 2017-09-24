package bootstrap3model

import (
	"github.com/adamcolton/buttress/bootstrap3"
)

type ValueBuilder func(string) string

type InputBuilder func(string, ValueBuilder) bootstrap3.Input

func HiddenTypeInput(name string, vg ValueBuilder) bootstrap3.Input {
	var v string
	if vg != nil {
		v = vg(name)
	}
	h := bootstrap3.NewHidden(name, v)
	return h
}

func TextTypeInput(name string, vg ValueBuilder) bootstrap3.Input {
	return Input("text", name, vg)
}

func NumberTypeInput(name string, vg ValueBuilder) bootstrap3.Input {
	return Input("number", name, vg)
}

func DatetimeLocalTypeInput(name string, vg ValueBuilder) bootstrap3.Input {
	return Input("datetime-local", name, vg)
}

func PasswordTypeInput(name string, vg ValueBuilder) bootstrap3.Input {
	return Input("password", name, vg)
}

func ConfirmPassword(name string, vg ValueBuilder) bootstrap3.Input {
	return bootstrap3.NewConfirmPassword(name)
}

func Input(kind, name string, vg ValueBuilder) *bootstrap3.InputTag {
	i := bootstrap3.NewInputTag(kind, name, name)
	if vg != nil {
		i.Value = vg(name)
	}
	return i
}

func GoTemplateValueBuilder(name string) string {
	return "{{." + name + "}}"
}

type TypeMap map[string]InputBuilder

var Types = TypeMap{
	//"bool":       gothicgo.BoolType,
	"byte":    NumberTypeInput,
	"int":     NumberTypeInput,
	"int8":    NumberTypeInput,
	"int16":   NumberTypeInput,
	"int32":   NumberTypeInput,
	"int64":   NumberTypeInput,
	"float32": NumberTypeInput,
	"float64": NumberTypeInput,
	//"rune":       gothicgo.RuneType,
	"string":   TextTypeInput,
	"uint":     NumberTypeInput,
	"uint8":    NumberTypeInput,
	"uint16":   NumberTypeInput,
	"uint32":   NumberTypeInput,
	"uint64":   NumberTypeInput,
	"datetime": DatetimeLocalTypeInput,
	"password": PasswordTypeInput,
}
