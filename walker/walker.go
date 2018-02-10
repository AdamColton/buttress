package walker

import (
	"os"
	"path/filepath"
	"regexp"
)

type DirFilter byte

const (
	DirAny DirFilter = iota
	DirYes
	DirNo
)

type Filter struct {
	Dir      DirFilter
	match    *regexp.Regexp
	Callback func(*Visit)
	del      []string
}

func New() *Filter {
	return &Filter{}
}

func (f *Filter) SetMatch(expr string) error {
	r, err := regexp.Compile(expr)
	if err != nil {
		return err
	}
	f.match = r
	return nil
}

type Visit struct {
	Path   string
	Info   os.FileInfo
	Err    error
	Delete bool
	// Matches holds the sub-matches from the regex
	Matches []string
}

func (f *Filter) walker(path string, info os.FileInfo, err error) error {
	if is := info.IsDir(); f.Dir != DirAny && (f.Dir == DirYes && !is || f.Dir == DirNo && is) {
		return nil
	}
	var m []string
	if f.match != nil {
		m = f.match.FindStringSubmatch(info.Name())
		if len(m) == 0 {
			return nil
		}
	}

	v := &Visit{
		Path:    path,
		Info:    info,
		Err:     err,
		Matches: m,
	}
	f.Callback(v)

	if v.Delete {
		f.del = append(f.del, path)
	}

	return v.Err
}

func (f *Filter) Walk(path string) {
	f.del = nil
	filepath.Walk(path, f.walker)
	for _, d := range f.del {
		os.Remove(d)
	}
}
