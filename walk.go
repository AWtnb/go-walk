package walk

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/AWtnb/go-everything"
)

const SEP = string(os.PathSeparator)

func getDepth(path string) int {
	return strings.Count(strings.TrimSuffix(path, SEP), SEP)
}

func isSubPath(base, target string) bool {
	base = filepath.Clean(base)
	target = filepath.Clean(target)

	rel, err := filepath.Rel(base, target)
	if err != nil {
		return false
	}
	return rel != "." && !strings.HasPrefix(rel, "..")
}

type Walker struct {
	all       bool
	root      string
	step      int
	exception []string
}

// `exclude` is comma-separated string
func (w *Walker) Init(root string, all bool, depth int, exclude string) {
	w.root = root
	w.all = all
	w.step = depth

	for _, e := range strings.Split(exclude, ",") {
		w.exception = append(w.exception, strings.TrimSpace(e))
	}
	a := "AppData"
	if !slices.Contains(w.exception, a) {
		w.exception = append(w.exception, a)
	}
}

func (w Walker) EverythingTraverse() (found []string, err error) {
	result, err := everything.Scan(w.root, !w.all)
	if err != nil {
		return
	}
	slices.Sort(result)

	// empty struct{} is 0 bytes in size, making it highly efficient
	skipPaths := make(map[string]struct{})
	rd := getDepth(w.root)

	for _, path := range result {

		skippable := false
		for sd := range skipPaths {
			if isSubPath(sd, path) {
				skippable = true
				break
			}
		}
		if skippable {
			continue
		}

		d := getDepth(path) - rd
		if d == 0 { // `path` is root itself
			found = append(found, path)
			continue
		}
		if -1 < w.step && w.step < d {
			d := filepath.Dir(path)
			skipPaths[d] = struct{}{}
			continue
		}

		name := filepath.Base(path)
		if slices.Contains(w.exception, name) || strings.HasPrefix(name, ".") {
			skipPaths[path] = struct{}{}
			continue
		}
		found = append(found, path)
	}
	return
}

func (w Walker) Traverse() (found []string, err error) {
	rd := getDepth(w.root)
	err = filepath.WalkDir(w.root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		d := getDepth(path) - rd
		if d == 0 { // `path` is root itself
			found = append(found, path)
			return nil
		}
		if -1 < w.step && w.step < d {
			return filepath.SkipDir
		}
		name := info.Name()
		if slices.Contains(w.exception, name) || strings.HasPrefix(name, ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() || w.all {
			found = append(found, path)
		}
		return nil
	})
	return
}
