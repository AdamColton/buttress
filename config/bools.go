package config

type BoolSetter interface {
	As(bool, ...string) boolSetter
}

type boolSetter string

func (s boolSetter) key() key {
	return key{
		name: string(s),
		kind: boolKind,
		env:  activeEnv,
	}
}

func (s boolSetter) As(val bool, environments ...string) boolSetter {
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

func SetBool(name string) boolSetter { return boolSetter(name) }
func GetBool(name string) (bool, error) {
	k := boolSetter(name).key()
	v, ok := vals[k]
	if !ok {
		return false, k.notFound()
	}
	return v.(bool), nil
}

func MustGetBool(name string) bool {
	s, err := GetBool(name)
	if err != nil {
		panic(err)
	}
	return s
}
