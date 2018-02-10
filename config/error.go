package config

import (
	"strings"
)

type ErrNotFound struct {
	Type, Env, Name string
}

func (e ErrNotFound) Error() string {
	msg := []string{"Cound not find ", e.Type, " '", e.Name, "' in environment '", e.Env, "'"}
	return strings.Join(msg, "")
}

var kindMap = map[int]string{
	strKind:   "string",
	boolKind:  "bool",
	bytesKind: "[]bytes",
}

func (k key) notFound() ErrNotFound {
	return ErrNotFound{
		Type: kindMap[k.kind],
		Env:  k.env,
		Name: k.name,
	}
}
