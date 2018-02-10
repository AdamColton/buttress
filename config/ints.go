package config

type IntSetter interface {
	As(int, ...string) intSetter
}

type intSetter string

func (s intSetter) key() key {
	return key{
		name: string(s),
		kind: intKind,
		env:  activeEnv,
	}
}

func (s intSetter) As(val int, environments ...string) intSetter {
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

func SetInt(name string) intSetter { return intSetter(name) }
func GetInt(name string) (int, error) {
	k := intSetter(name).key()
	v, ok := vals[k]
	if !ok {
		return 0, k.notFound()
	}
	return v.(int), nil
}

func MustGetInt(name string) int {
	s, err := GetInt(name)
	if err != nil {
		panic(err)
	}
	return s
}
