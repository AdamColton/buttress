package csscontext

type CSSContext interface {
	String() string
	private()
}

type csscontext string

func (c csscontext) String() string { return string(c) }

func (csscontext) private() {}

func Default() CSSContext { return csscontext("default") }
func Primary() CSSContext { return csscontext("primary") }
func Success() CSSContext { return csscontext("success") }
func Info() CSSContext    { return csscontext("info") }
func Warning() CSSContext { return csscontext("warning") }
func Danger() CSSContext  { return csscontext("danger") }
