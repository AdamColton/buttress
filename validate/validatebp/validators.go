package validatebp

import (
	"github.com/adamcolton/gothic/gothicgo"
	"io"
)

type StructWriterTo func(value, name string, strct *gothicgo.Struct) io.WriterTo

type FieldWriterTo func(value, name string, field *gothicgo.Field, strct *gothicgo.Struct) io.WriterTo

var structValidators = make(map[string]StructWriterTo)
var fieldValidators = make(map[string]FieldWriterTo)

const ErrAlreadDefined errStr = "Validator already defined"
const ErrCannotBeNil errStr = "Cannot be nil"

func RegisterStructWriterTo(name string, validator StructWriterTo) error {
	if _, ok := structValidators[name]; ok {
		return ErrAlreadDefined
	}
	if validator == nil {
		return ErrCannotBeNil
	}
	structValidators[name] = validator
	return nil
}

func RegisterFieldWriterTo(name string, validator FieldWriterTo) error {
	if _, ok := fieldValidators[name]; ok {
		return ErrAlreadDefined
	}
	if validator == nil {
		return ErrCannotBeNil
	}
	fieldValidators[name] = validator
	return nil
}
