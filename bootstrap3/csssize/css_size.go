package csssize

type CSSSize interface {
	String() string
	Smaller() CSSSize
	Bigger() CSSSize
	privateCSSSize()
}

type csssize string

func (s csssize) String() string { return string(s) }
func (csssize) privateCSSSize()  {}

func (s csssize) Smaller() CSSSize {
	switch s.String() {
	case "sm":
		return ExtraSmall()
	case "md":
		return Small()
	case "lg":
		return Medium()
	}
	return nil
}

func (s csssize) Bigger() CSSSize {
	switch s.String() {
	case "xs":
		return Small()
	case "sm":
		return Medium()
	case "lg":
		return Large()
	}
	return nil
}

func ExtraSmall() CSSSize { return csssize("xs") }
func Small() CSSSize      { return csssize("sm") }
func Medium() CSSSize     { return csssize("md") }
func Large() CSSSize      { return csssize("lg") }

func Sizes() []CSSSize {
	return []CSSSize{
		ExtraSmall(),
		Small(),
		Medium(),
		Large(),
	}
}
