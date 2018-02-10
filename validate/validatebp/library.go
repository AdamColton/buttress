package validatebp

import (
	"fmt"
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicio"
	"io"
)

var tooShort = `if len(%s.%s) < %s && len(%s.%s) > 0 {
	r.AddToField("%s", "Too Short")
}
`

func MinLength(value, name string, field *gothicgo.Field, strct *gothicgo.Struct) io.WriterTo {
	s := fmt.Sprintf(tooShort, strct.ReceiverName, field.Name(), value, strct.ReceiverName, field.Name(), field.Name())
	return gothicio.StringWriterTo(s)
}

var tooLong = `if len(%s.%s) > %s {
	r.AddToField("%s", "Too Long")
}
`

func MaxLength(value, name string, field *gothicgo.Field, strct *gothicgo.Struct) io.WriterTo {
	s := fmt.Sprintf(tooLong, strct.ReceiverName, field.Name(), value, field.Name())
	return gothicio.StringWriterTo(s)
}

var requiredByLen = `if len(%s.%s) == 0 {
	r.AddToField("%s", "Required")
}
`

var requiredByPtr = `if %s.%s == nil {
	r.AddToField("%s", "Required")
}
`

func Required(value, name string, field *gothicgo.Field, strct *gothicgo.Struct) io.WriterTo {
	sf, ok := strct.Field(field.Name())
	if !ok {
		return nil
	}
	var str string
	if _, ok := sf.Type().(gothicgo.PointerType); ok {
		str = fmt.Sprintf(requiredByPtr, strct.ReceiverName, field.Name(), field.Name())
	} else {
		str = fmt.Sprintf(requiredByLen, strct.ReceiverName, field.Name(), field.Name())
	}
	return gothicio.StringWriterTo(str)
}

var tooLow = `if %s.%s < %s {
	r.AddToField("%s", "Too Low")
}
`

func MinRange(value, name string, field *gothicgo.Field, strct *gothicgo.Struct) io.WriterTo {
	s := fmt.Sprintf(tooLow, strct.ReceiverName, field.Name(), value, field.Name())
	return gothicio.StringWriterTo(s)
}

var tooHigh = `if %s.%s > %s {
	r.AddToField("%s", "Too High")
}
`

func MaxRange(value, name string, field *gothicgo.Field, strct *gothicgo.Struct) io.WriterTo {
	s := fmt.Sprintf(tooHigh, strct.ReceiverName, field.Name(), value, field.Name())
	return gothicio.StringWriterTo(s)
}
