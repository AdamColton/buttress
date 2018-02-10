package config

import (
	"encoding/base64"
)

var DecodeString = base64.URLEncoding.DecodeString
var EncodeString = base64.URLEncoding.EncodeToString

type BytesSetter interface {
	As(val []byte, environments ...string) BytesSetter
	AsBase64(val string, environments ...string) (BytesSetter, error)
	MustAsBase64(val string, environments ...string) BytesSetter
}

type bytesSetter string

func (s bytesSetter) key() key {
	return key{
		name: string(s),
		kind: bytesKind,
		env:  activeEnv,
	}
}

func (s bytesSetter) As(val []byte, environments ...string) BytesSetter {
	if len(environments) == 0 {
		environments = allenvironments
	}

	k := s.key()
	for _, envStr := range environments {
		k.env = envStr
		vals[k] = val
	}
	return s
}

func (s bytesSetter) AsBase64(val string, environments ...string) (BytesSetter, error) {
	b, err := DecodeString(val)
	if err != nil {
		return nil, err
	}
	return s.As(b, environments...), nil
}

func (s bytesSetter) MustAsBase64(val string, environments ...string) BytesSetter {
	bs, err := s.AsBase64(val, environments...)
	if err != nil {
		panic(err)
	}
	return bs
}

func SetBytes(name string) BytesSetter { return bytesSetter(name) }
func GetBytes(name string) ([]byte, error) {
	k := bytesSetter(name).key()
	v, ok := vals[k]
	if !ok {
		return nil, k.notFound()
	}
	return v.([]byte), nil
}

func MustGetBytes(name string) []byte {
	s, err := GetBytes(name)
	if err != nil {
		panic(err)
	}
	return s
}
