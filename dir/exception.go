package dir

import (
	"os"
	"strings"
)

type WalkException struct {
	names []string
}

func (wex *WalkException) SetName(s string) {
	wex.names = append(wex.names, strings.TrimSpace(s))
}

func (wex *WalkException) SetNames(s string, sep string) {
	if len(s) < 1 {
		return
	}
	for _, elem := range strings.Split(s, sep) {
		wex.SetName(elem)
	}
}

func (wex WalkException) Contains(name string) bool {
	for _, n := range wex.names {
		if n == name {
			return true
		}
	}
	return false
}

func (wex WalkException) isSkippablePath(path string) bool {
	sep := string(os.PathSeparator)
	if strings.HasPrefix(path, ".") || strings.Contains(path, sep+".") {
		return true
	}
	for _, n := range wex.names {
		if strings.HasPrefix(path, n) || strings.Contains(path, sep+n+sep) || strings.HasSuffix(path, n) {
			return true
		}
	}
	return false
}

func (wex WalkException) Filter(paths []string, root string) (filtered []string) {
	if len(wex.names) < 1 {
		filtered = paths
		return
	}
	for i := 0; i < len(paths); i++ {
		p := paths[i]
		rel := strings.TrimPrefix(strings.TrimPrefix(p, root), string(os.PathSeparator))
		if wex.isSkippablePath(rel) {
			continue
		}
		filtered = append(filtered, p)
	}
	return
}
