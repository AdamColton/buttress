package gothicproject

import (
	"os"
	"path/filepath"
	"time"
)

func GoFiles(dir string) ([]string, error) {
	return filepath.Glob(filepath.Join(dir, "*.go"))
}

// MostRecentModTime returns the most recently modfied file in a directory
func MostRecentModTime(dir string) time.Time {
	mrmtw := &mostRecentModeTimeWalker{}
	filepath.Walk(dir, mrmtw.walkDir)
	return mrmtw.Time
}

type mostRecentModeTimeWalker struct {
	time.Time
}

func (mrmtw *mostRecentModeTimeWalker) walkDir(path string, info os.FileInfo, err error) error {
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
