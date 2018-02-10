package validatebp

// technically, more of a lexer; parse takes a string and returns it as a list
// of validators.
func Parse(str string, vals map[string]string) map[string]string {
	var key string
	var val, spc []rune
	endOfKey := ':'
	rs := []rune(str)
	for i := 0; i < len(rs); i++ {
		switch c := rs[i]; c {
		case endOfKey:
			key = string(val)
			endOfKey = 0
			val = nil
		case '|':
			if endOfKey == ':' {
				vals[string(val)] = ""
			} else {
				vals[key] = string(val)
				endOfKey = ':'
			}
			spc = nil
			val = nil
		case ' ', '\t':
			if len(val) > 0 {
				spc = append(spc, c)
			}
		case '\\':
			if i+1 < len(rs) {
				i++
				c = rs[i]
			} else {
				break
			}
			fallthrough
		default:
			val = append(val, spc...)
			spc = nil
			val = append(val, c)
		}
	}

	if endOfKey == ':' {
		vals[string(val)] = ""
	} else {
		vals[key] = string(val)
	}
	return vals
}

type errStr string

func (e errStr) Error() string {
	return string(e)
}
