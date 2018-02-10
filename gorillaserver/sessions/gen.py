typeFn = """
func (s *Session) <<Type>>(name <<string>>, dflt <<type>>) <<type>> {
  if ifc, ok := s.Values[name]; ok {
    if val, ok := ifc.(<<type>>); ok {
      return val
    }
  }
  return dflt
}
"""

types = [
  "string",
  "int",
  "uint",
  "uint64",
  "float32",
  "float64",
  ("[]byte","Bytes"),
]

f = open("types.go", "w")
f.write("package sessions\n")

for t in types:
  if type(t)==str:
    l = t
    u = t[0].upper()+t[1:]
  else:
    l=t[0]
    u=t[1]
  fn = typeFn.replace("<<type>>",l).replace("<<Type>>",u)
  if l=="string":
    fn = fn.replace("<<string>>","")
  else:
    fn = fn.replace("<<string>>","string")
  f.write(fn)