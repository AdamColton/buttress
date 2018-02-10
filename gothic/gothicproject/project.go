package gothicproject

import (
	"encoding/json"
	"fmt"
	"github.com/adamcolton/buttress/walker"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	TimeFormat = "200601021504"
	App        = "app"
	Blueprint  = "blueprint"
)

type GothicProject struct {
	Name       string
	App        string
	Blueprint  string
	Tests      []string
	TimeFormat string
	Watch      []string
	gi         *GothicInstance
}

const gothicInstanceFilename = "gothic.instance.json"

type GothicInstance struct {
	LastGenerate     time.Time
	lastBlueprintMod time.Time
}

func LoadProject(r io.Reader) (*GothicProject, error) {
	var gp GothicProject
	err := json.NewDecoder(r).Decode(&gp)
	if err != nil {
		return nil, err
	}
	gp.App = SetDefault(gp.App, App)
	gp.Blueprint = SetDefault(gp.Blueprint, Blueprint)
	gp.TimeFormat = SetDefault(gp.TimeFormat, TimeFormat)
	if gp.Watch == nil {
		gp.Watch = []string{"go"}
	}
	return &gp, nil
}

func LoadInstance(r io.Reader) (*GothicInstance, error) {
	var gi GothicInstance
	err := json.NewDecoder(r).Decode(&gi)
	if err != nil {
		return nil, err
	}
	return &gi, nil
}

func (gi *GothicInstance) Save() error {
	f, err := os.Create(gothicInstanceFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(gi)
}

func findGothicFile() *os.File {
	prev, err := os.Getwd()
	if err != nil {
		return nil
	}
	restore := prev
	for err == nil {
		f, err := os.Open("gothic.json")
		if err == nil && f != nil {
			return f
		}
		err = os.Chdir("..")
		if err != nil {
			break
		}
		d, err := os.Getwd()
		if err != nil || d == prev {
			break
		}
		prev = d
	}
	os.Chdir(restore)
	return nil
}

// Load tries to load the Gothic Project. It will look up the directory tree for
// gothic.json. If it finds it, it will leave the working directory there and
// return the project.
func Load() (*GothicProject, error) {
	gf := findGothicFile()
	if gf == nil {
		return nil, fmt.Errorf("Not in a Gothic project directory")
	}
	defer gf.Close()
	gp, err := LoadProject(gf)
	if err != nil {
		return nil, err
	}

	gif, _ := os.Open(gothicInstanceFilename)
	if gif != nil {
		gp.gi, err = LoadInstance(gif)
		if err != nil {
			return nil, err
		}
	} else {
		gp.gi = &GothicInstance{}
	}

	return gp, nil
}

var extRe = regexp.MustCompile(`^\w+$`)

func (gp *GothicProject) watchToRe() (string, error) {
	if len(gp.Watch) == 0 {
		return "", fmt.Errorf("No file extensions to watch")
	}
	for _, ext := range gp.Watch {
		if !extRe.MatchString(ext) {
			return "", fmt.Errorf("File extensions to watch should be alpha-numeric and underscore: %s", ext)
		}
	}
	return "(" + strings.Join(gp.Watch, ")|(") + ")", nil
}

func (gp *GothicProject) RunBlueprint(force bool) (string, error) {
	if !force {
		reStr, err := gp.watchToRe()
		if err != nil {
			return "", err
		}
		w := walker.New()
		if err := w.SetMatch(reStr); err != nil {
			return "", err
		}
		w.Dir = walker.DirNo
		w.Callback = gp.gi.checkModTime
		w.Walk(gp.Blueprint)

		if !gp.gi.lastBlueprintMod.After(gp.gi.LastGenerate) {
			return "", nil
		}
	}
	fmt.Println("Generating")

	files, err := GoFiles(gp.Blueprint)
	if err != nil {
		return "", err
	}
	args := append([]string{"run"}, files...)
	out, err := exec.Command("go", args...).CombinedOutput()
	if err == nil {
		gp.gi.LastGenerate = time.Now()
		gp.gi.Save()
	}

	return string(out), err
}

func (gp *GothicProject) Instance(force bool) *GothicInstance {
	return gp.gi
}

func (gi *GothicInstance) checkModTime(v *walker.Visit) {
	if t := v.Info.ModTime(); gi.lastBlueprintMod.Before(t) {
		gi.lastBlueprintMod = t
	}
}
