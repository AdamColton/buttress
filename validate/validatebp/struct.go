package validatebp

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicio"
	"io"
)

// StructValidator can add validators methods to a struct
type StructValidator struct {
	Struct          *gothicgo.Struct
	validators      map[string]string
	fieldValidators map[string]map[string]string
}

// Struct returns a new StructValidator for a Struct.
func Struct(strct *gothicgo.Struct) *StructValidator {
	return &StructValidator{
		Struct:          strct,
		validators:      make(map[string]string),
		fieldValidators: make(map[string]map[string]string),
	}
}

func (s *StructValidator) Add(str string) {
	Parse(str, s.validators)
}

func (s *StructValidator) AddField(name, str string) {
	fieldValidators, ok := s.fieldValidators[name]
	if !ok {
		fieldValidators = make(map[string]string)
		s.fieldValidators[name] = fieldValidators
	}
	Parse(str, fieldValidators)
}

var pkg = gothicgo.MustPackageRef("github.com/adamcolton/buttress/validate")
var valResultType = gothicgo.PointerTo(gothicgo.DefStruct(pkg, "Results"))

func (s *StructValidator) Validate(validateStruct bool, fields ...string) io.WriterTo {
	var writerTo gothicio.WriterTos
	writerTo = append(writerTo, newResultsWriterTo{s.Struct.File()})
	if validateStruct {
		for k, v := range s.validators {
			if sv, ok := structValidators[k]; ok {
				writerTo = append(writerTo, sv(v, s.Struct.ReceiverName, s.Struct))
			}
		}
	}

	for _, fn := range fields {
		f, ok := s.Struct.Field(fn)
		if !ok {
			continue
		}
		for k, v := range s.fieldValidators[fn] {
			if fv, ok := fieldValidators[k]; ok {
				writerTo = append(writerTo, fv(v, s.Struct.ReceiverName, f, s.Struct))
			}
		}
	}
	return writerTo
}

func (s *StructValidator) ValidateMethod(name string, validateStruct bool, fields ...string) *gothicgo.Method {
	fn := s.Struct.NewMethod(name)
	s.Struct.File().AddRefImports(pkg)
	fn.UnnamedReturns(valResultType)
	fn.Body = gothicio.WriterToMerge(s.Validate(validateStruct, fields...), gothicio.StringWriterTo("return r\n"))

	return fn
}

type newResultsWriterTo struct {
	f *gothicgo.File
}

var caller = gothicgo.FuncCall(pkg, "New", nil, gothicgo.Rets(valResultType))

func (r newResultsWriterTo) WriteTo(w io.Writer) (int64, error) {
	s := gothicio.NewSumWriter(w)
	s.WriteString("r := ")
	s.WriteString(caller.Call(r.f))
	s.WriteString("\n")
	return s.Sum, s.Err
}
