package config

type StrSetter interface {
	As(string, ...string) StrSetter
}

type strSetter string

func (s strSetter) key() key {
	return key{
		name: string(s),
		kind: strKind,
		env:  activeEnv,
	}
}

func (s strSetter) As(val string, environments ...string) StrSetter {
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

func SetString(name string) StrSetter { return strSetter(name) }
func GetString(name string) (string, error) {
	k := strSetter(name).key()
	v, ok := vals[k]
	if !ok {
		return "", k.notFound()
	}
	return v.(string), nil
}

func MustGetString(name string) string {
	s, err := GetString(name)
	if err != nil {
		panic(err)
	}
	return s
}
