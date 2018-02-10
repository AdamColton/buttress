package config

type FloatSetter interface {
	As(float64, ...string) floatSetter
}

type floatSetter string

func (s floatSetter) key() key {
	return key{
		name: string(s),
		kind: floatKind,
		env:  activeEnv,
	}
}

func (s floatSetter) As(val float64, environments ...string) floatSetter {
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

func SetFloat(name string) floatSetter { return floatSetter(name) }
func GetFloat(name string) (float64, error) {
	k := floatSetter(name).key()
	v, ok := vals[k]
	if !ok {
		return 0, k.notFound()
	}
	return v.(float64), nil
}

func MustGetFloat(name string) float64 {
	s, err := GetFloat(name)
	if err != nil {
		panic(err)
	}
	return s
}
