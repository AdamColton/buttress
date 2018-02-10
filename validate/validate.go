package validate

type Results struct {
	Object []string
	Field  map[string][]string
	OK     bool
}

func New() *Results {
	return &Results{
		Field: make(map[string][]string),
		OK:    true,
	}
}

func (r *Results) AddToField(name, message string) {
	r.OK = false
	r.Field[name] = append(r.Field[name], message)
}
