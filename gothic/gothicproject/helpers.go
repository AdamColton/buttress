package gothicproject

import (
	"os"
	"path/filepath"
	"time"
)

func GoFiles(dir string) ([]string, error) {
	return filepath.Glob(filepath.Join(dir, "*.go"))
}

func (mrmtw *mostRecentModTimeWalker) walkDir(path string, info os.FileInfo, err error) error {
	if t := info.ModTime(); mrmtw.Before(t) {
		mrmtw.Time = t
	}
	return nil
}

func SetDefault(cur, dflt string) string {
	if cur == "" {
		return dflt
	}
	return cur
}

type mostRecentModTimeWalker struct {
	time.Time
}
